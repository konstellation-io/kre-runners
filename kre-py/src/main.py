import asyncio
import copy
import importlib.util
import inspect
import os
import sys
import traceback
import gc
from datetime import datetime

from nats.js.api import DeliverPolicy, ConsumerConfig

from compression import compress_if_needed, is_compressed, uncompress
from config import Config
from context import HandlerContext
from kre_nats_msg_pb2 import KreNatsMessage
from kre_runner import Runner

from nats.errors import ConnectionClosedError, TimeoutError


class NodeRunner(Runner):
    def __init__(self):
        config = Config()
        name = f"{config.krt_version}-{config.krt_node_name}"
        Runner.__init__(self, name, config)
        self.handler_ctx = None
        self.handler_init_fn = None
        self.handler_fn = None
        self.load_handler()

    def load_handler(self) -> None:
        self.logger.info(f"Loading handler script {self.config.handler_path}...")

        handler_full_path = os.path.join(self.config.base_path, self.config.handler_path)
        handler_dirname = os.path.dirname(handler_full_path)
        sys.path.append(handler_dirname)

        try:
            spec = importlib.util.spec_from_file_location("worker", handler_full_path)
            handler_module = importlib.util.module_from_spec(spec)
            spec.loader.exec_module(handler_module)
        except Exception as err:
            tb = traceback.format_exc()
            self.logger.error(
                f"Error loading handler script {self.config.handler_path}: {err}\n\n{tb}"
            )
            sys.exit(1)

        if not hasattr(handler_module, "handler"):
            raise Exception(
                f"Handler module '{handler_full_path}' must "
                f"implement a function 'handler(ctx, data)'"
            )

        if hasattr(handler_module, "init"):
            self.handler_init_fn = handler_module.init

        self.handler_fn = handler_module.handler
        self.logger.info(f"Handler script was loaded from '{handler_full_path}'")

    async def execute_handler_init(self) -> None:
        self.logger.info(f"Creating handler context...")
        self.handler_ctx = HandlerContext(
            self.config,
            self.nc,
            self.mongo_conn,
            self.logger,
            self.early_reply,
        )

        if not self.handler_init_fn:
            return

        self.logger.info(f"Executing handler init...")

        if inspect.iscoroutinefunction(self.handler_init_fn):
            await asyncio.create_task(self.handler_init_fn(self.handler_ctx))
        else:
            self.handler_init_fn(self.handler_ctx)

    async def process_messages(self) -> None:
        queue_name = f"queue_{self.config.nats_input}"

        try:
            self.subscription_sid = await self.js.subscribe(
                stream=self.config.nats_stream,
                subject=self.config.nats_input,
                queue=self.runner_name,
                durable=self.runner_name,
                cb=self.create_message_cb(),
                config=ConsumerConfig(
                    deliver_policy=DeliverPolicy.NEW,
                    ack_wait = 22 * 3600,  # 22 hours
                ),
                manual_ack=True,
            )
        except Exception as err:
            tb = traceback.format_exc()
            self.logger.error(
                f"Error subscribing to stream {self.config.handler_path}: {err}\n\n{tb}"
            )
            sys.exit(1)

        self.logger.info(
            f"Listening to '{self.config.nats_input}' subject with queue '{queue_name}' "
            f"from stream '{self.config.nats_stream}'"
        )

        await self.execute_handler_init()

    def create_message_cb(self) -> callable:
        async def message_cb(msg) -> None:

            """
            Callback for processing a message.

            :param msg: The message to process.
            """

            start = datetime.utcnow()

            # Parse incoming message
            request_msg = self.new_request_msg(msg.data)

            self.logger.info(f"Received new request message from NATS subject {msg.subject}")

            try:
                # Make a shallow copy of the ctx object to set inside the request msg.
                ctx = copy.copy(self.handler_ctx)
                ctx.__request_msg__ = request_msg

                # Execute the handler function sending context and the payload.
                handler_result = await self.handler_fn(ctx, request_msg.payload)
                # Tell NATS we don't need to receive the message anymore and we are done processing it.
                await msg.ack()
                end = datetime.utcnow()

                # Save the elapsed time for this node and for the workflow if it is the last node.
                self.save_elapsed_time(request_msg, start, end)

                # Ignore send reply if the msg was replied previously.
                if self.config.krt_last_node == "true" and request_msg.replied:
                    if handler_result is not None:
                        self.logger.info(
                            "ignoring the last node response because the message was replied previously"
                        )
                    return

                # Generate a KreNatsMessage response.
                res = self.new_response_msg(request_msg, handler_result, start, end)

                # Publish the response message to the output subject.
                output_subject = self.get_output_subject(request_msg.early_exit)
                await self.publish_response(output_subject, res)

            except Exception as err:
                # Publish an error message to the final reply subject
                # in order to stop the workflow execution. So the next nodes will be ignored
                # and the gRPC response will be an exception.

                # Tell NATS we don't need to receive the message anymore and we are done processing it.
                await msg.ack()

                tb = traceback.format_exc()
                self.logger.error(f"error executing handler: {err} \n\n{tb}")
                response_err = KreNatsMessage()
                response_err.error = f"error in '{self.config.krt_node_name}': {str(err)}"
                output_subject = self.config.nats_entrypoint_subject

                await self.js.publish(
                    stream=self.config.nats_stream,
                    subject=output_subject,
                    payload=response_err.SerializeToString(),
                )
            gc.collect()

        return message_cb

    def get_output_subject(self, early_exit: bool) -> str:
        if early_exit:
            self.logger.info(f"Early exit recieved, worklow has stopped execution")
            output_subject = self.config.nats_entrypoint_subject
        else:
            output_subject = self.config.nats_output
        return output_subject

    @staticmethod
    def new_request_msg(data: bytes) -> KreNatsMessage:
        if is_compressed(data):
            data = uncompress(data)

        request_msg = KreNatsMessage()
        request_msg.ParseFromString(data)

        return request_msg

    # new_response_msg creates a KreNatsMessage maintaining the tracking ID plus adding the
    # handler result and the tracking information for this node.
    def new_response_msg(
            self, request_msg: KreNatsMessage, payload: any, start, end
    ) -> KreNatsMessage:
        res = KreNatsMessage()
        res.replied = request_msg.replied
        res.request_id = request_msg.request_id
        res.tracking_id = request_msg.tracking_id
        res.tracking.extend(request_msg.tracking)
        res.payload.Pack(payload)

        # add tracking info
        t = res.tracking.add()
        t.node_name = self.config.krt_node_name
        t.start = start.isoformat()
        t.end = end.isoformat()

        return res

    async def publish_response(self, subject: str, response_msg: KreNatsMessage) -> None:

        """
        Publish the response message to the output subject.

        Args:
            subject: The subject to publish the response message.
            response_msg:  The response message to publish.
        """

        serialized_response_msg = compress_if_needed(
            response_msg.SerializeToString(), logger=self.logger
        )

        try:
            await self.js.publish(
                stream=self.config.nats_stream, subject=subject, payload=serialized_response_msg
            )
            self.logger.info(
                f"Published response to NATS subject {subject} "
                f"from stream {self.config.nats_stream}"
            )
        except ConnectionClosedError as err:
            self.logger.error(f"Connection closed when publishing response: {err}")
        except TimeoutError as err:
            self.logger.error(f"Timeout when publishing response: {err}")
        except Exception as err:
            self.logger.error(f"Error publishing response: {err}")

    async def early_reply(self, response: any, request_id: str):
        res = KreNatsMessage()
        res.payload.Pack(response)
        res.request_id = request_id
        await self.publish_response(self.config.nats_entrypoint_subject, res)

    def save_elapsed_time(self, req_msg: KreNatsMessage, start: datetime, end: datetime) -> None:
        """
        save_elapsed_time stores in InfluxDB the elapsed time for the current node and the
        total elapsed time of the complete workflow if it is the last node.

        :param req_msg: the request message.
        :param start: when this node started.
        :param end: when this node ended.
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

        is_last_node = True if self.config.krt_last_node == "true" else False

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
