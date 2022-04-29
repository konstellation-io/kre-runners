import asyncio
import copy
import importlib.util
import inspect
import os
import sys
import traceback
from datetime import datetime

import pymongo

from compression import compress_if_needed, is_compressed, uncompress
from config import Config
from context import HandlerContext
from kre_nats_msg_pb2 import KreNatsMessage
from kre_runner import Runner


class NodeRunner(Runner):
    def __init__(self):
        config = Config()
        name = f"{config.krt_version}-{config.krt_node_name}"
        Runner.__init__(self, name, config)
        self.handler_ctx = None
        self.handler_init_fn = None
        self.handler_fn = None
        self.mongo_conn = None
        self.load_handler()

    def load_handler(self):
        self.logger.info(f"loading handler script {self.config.handler_path}...")

        handler_full_path = os.path.join(
            self.config.base_path, self.config.handler_path
        )
        handler_dirname = os.path.dirname(handler_full_path)
        sys.path.append(handler_dirname)

        try:
            spec = importlib.util.spec_from_file_location("worker", handler_full_path)
            handler_module = importlib.util.module_from_spec(spec)
            spec.loader.exec_module(handler_module)
        except Exception as err:
            tb = traceback.format_exc()
            self.logger.error(
                f"error loading handler script {self.config.handler_path}: {err}\n\n{tb}"
            )
            sys.exit(1)

        if not hasattr(handler_module, "handler"):
            raise Exception(
                f"handler module '{handler_full_path}' must implement a function 'handler(ctx, data)'"
            )

        if hasattr(handler_module, "init"):
            self.handler_init_fn = handler_module.init

        self.handler_fn = handler_module.handler
        self.logger.info(f"handler script was loaded from '{handler_full_path}'")

    async def execute_handler_init(self):
        self.logger.info(f"creating handler context...")
        self.handler_ctx = HandlerContext(
            self.config,
            self.nc,
            self.mongo_conn,
            self.logger,
            self.early_reply,
            self.send_output,
        )

        if not self.handler_init_fn:
            return

        self.logger.info(f"executing handler init...")
        if inspect.iscoroutinefunction(self.handler_init_fn):
            await asyncio.create_task(self.handler_init_fn(self.handler_ctx))
        else:
            self.handler_init_fn(self.handler_ctx)

    async def process_messages(self):
        self.logger.info(f"connecting to MongoDB...")
        self.mongo_conn = pymongo.MongoClient(
            self.config.mongo_uri, socketTimeoutMS=10000, connectTimeoutMS=10000
        )

        queue_name = f"queue_{self.config.nats_input}"
        self.subscription_sid = await self.nc.subscribe(
            self.config.nats_input, cb=self.create_message_cb(), queue=queue_name
        )
        self.logger.info(
            f"listening to '{self.config.nats_input}' subject with queue '{queue_name}'"
        )

        await self.execute_handler_init()

    def create_message_cb(self):
        async def message_cb(msg):
            start = datetime.utcnow()

            # Parse incoming message
            request_msg = self.new_request_msg(msg.data)

            try:
                # Validations
                if msg.reply == "" and request_msg.reply == "":
                    raise Exception("the reply subject was not found")

                # If the "msg.reply" has a value, it means that is the first message of the workflow.
                # In other words, it comes from the initial entrypoint request.
                # In this case we set this value into the "request_msg.reply" field in order
                # to be propagated for the rest of workflow nodes.
                if msg.reply != "":
                    request_msg.reply = msg.reply

                self.logger.info(
                    f"received message on '{msg.subject}' with final reply '{request_msg.reply}'"
                )

                # Make a shallow copy of the ctx object to set inside the request msg.
                ctx = copy.copy(self.handler_ctx)
                ctx.__request_msg__ = request_msg

                # Execute the handler function sending context and the payload.
                handler_result = await self.handler_fn(ctx, request_msg.payload)
                end = datetime.utcnow()

                # Save the elapsed time for this node and for the workflow if it is the last node.
                is_last_node = self.config.nats_output == ""
                if not request_msg.isIntermediateMessage:
                    self.save_elapsed_time(request_msg, start, end, is_last_node)

                # Ignore send reply if the msg was replied previousl y.
                if is_last_node and request_msg.replied:
                    if handler_result is not None:
                        self.logger.info(
                            "ignoring the last node response because the message was replied previously"
                        )
                    return

                # Generate a KreNatsMessage response.
                res = self.new_response_msg(request_msg, handler_result, start, end)

                # Publish the response message to the output subject.
                output_subject = (
                    request_msg.reply if is_last_node else self.config.nats_output
                )
                await self.publish_response(output_subject, res)

            except Exception as err:
                tb = traceback.format_exc()
                self.logger.error(f"error executing handler: {err} \n\n{tb}")
                response_err = KreNatsMessage()
                response_err.error = (
                    f"error in '{self.config.krt_node_name}': {str(err)}"
                )
                await self.nc.publish(
                    request_msg.reply, response_err.SerializeToString()
                )

        return message_cb

    @staticmethod
    def new_request_msg(data: bytes) -> KreNatsMessage:
        if is_compressed(data):
            data = uncompress(data)

        request_msg = KreNatsMessage()
        request_msg.ParseFromString(data)

        return request_msg

    # new_response_msg creates a KreNatsMessage maintaining the tracking ID and adding the
    # handler result and the tracking information for this node.
    def new_response_msg(
        self, request_msg: KreNatsMessage, payload: any, start, end
    ) -> KreNatsMessage:
        res = KreNatsMessage()
        res.replied = request_msg.replied
        res.reply = request_msg.reply
        res.tracking_id = request_msg.tracking_id
        res.tracking.extend(request_msg.tracking)
        res.payload.Pack(payload)

        # add tracking info
        t = res.tracking.add()
        t.node_name = self.config.krt_node_name
        t.start = start.isoformat()
        t.end = end.isoformat()

        return res

    async def publish_response(self, subject, response_msg):
        serialized_response_msg = compress_if_needed(
            response_msg.SerializeToString(), logger=self.logger
        )
        await self.nc.publish(subject, serialized_response_msg)

        self.logger.info(f"published response to NATS subject '{subject}'")
        await self.nc.flush(timeout=self.config.nats_flush_timeout)

    async def early_reply(self, nats_reply_subject: str, response: any):
        res = KreNatsMessage()
        res.payload.Pack(response)
        await self.publish_response(nats_reply_subject, res)

    async def send_output(self, nats_output: str, nats_reply_subject: str, response: any):
        res = KreNatsMessage()
        res.payload.Pack(response)
        res.isIntermediateMessage = True
        res.reply = nats_reply_subject
        await self.publish_response(nats_output, res)

    def save_elapsed_time(
        self,
        req_msg: KreNatsMessage,
        start: datetime,
        end: datetime,
        is_last_node: bool,
    ) -> None:
        """
        save_elapsed_time stores in InfluxDB the elapsed time for the current node and the total elapsed time
        of the complete workflow if it is the last node.

        :param req_msg: the request message.
        :param start: when this node started.
        :param end: when this node ended.
        :param is_last_node: indicates if this node is the last one.
        :return: None
        """
        # Save the elapsed time for this node
        prev = req_msg.tracking[-1]
        prev_end = datetime.fromisoformat(prev.end)

        elapsed = end - start
        waiting = start - prev_end

        tags = {
            "workflow": self.config.krt_workflow_name,
            "version": self.config.krt_version,
            "node": self.config.krt_node_name,
        }

        fields = {
            "tracking_id": req_msg.tracking_id,
            "node_from": prev.node_name,
            "elapsed_ms": elapsed.total_seconds() * 1000,
            "waiting_ms": waiting.total_seconds() * 1000,
        }

        self.handler_ctx.measurement.save("node_elapsed_time", fields, tags)

        if is_last_node:
            # Save the complete workflow elapsed time
            entrypoint = req_msg.tracking[0]
            entrypoint_start = datetime.fromisoformat(entrypoint.start)
            elapsed = (end - entrypoint_start).total_seconds() * 1000

            tags = {
                "workflow": self.config.krt_workflow_name,
                "version": self.config.krt_version,
            }
            fields = {"elapsed_ms": elapsed, "tracking_id": req_msg.tracking_id}
            self.handler_ctx.measurement.save("workflow_elapsed_time", fields, tags)


if __name__ == "__main__":
    runner = NodeRunner()
    runner.start()
