import abc
from datetime import datetime
import gzip
import traceback
import uuid
import asyncio

from grpclib.server import Stream
from grpclib import GRPCError
from grpclib.const import Status
from nats.js.api import ConsumerConfig, DeliverPolicy, StreamConfig, RetentionPolicy
from nats.aio.client import Client as NATS

from kre_nats_msg_pb2 import KreNatsMessage

COMPRESS_LEVEL = 9
MESSAGE_THRESHOLD = 1024 * 1024
GZIP_HEADER = b"\x1f\x8b"


# NOTE: EntrypointKRE will be extended by Entrypoint class auto-generated
class EntrypointKRE:
    def __init__(self, logger, subjects, config):
        self.logger = logger
        self.nc = NATS()
        self.js = None
        self.subscriptions = {}
        self.streams = {}
        self.grpc_streams = {}
        self.subjects = subjects
        self.config = config

    @abc.abstractmethod
    def make_response_object(self, subject: str, response: KreNatsMessage) -> bytes:
        # To be implemented on autogenerated entrypoint
        pass

    async def stop(self):
        """
        Stop the Entrypoint service by closing the NATS connection
        """
        if not self.nc.is_closed:
            self.logger.info("closing NATS connection")
            await self.nc.close()

    async def start(self) -> None:
        """
        Starts the entrypoint service by connecting to the NATS server and subscribing to the
        subjects related to each workflow exposed by the Entrypoint.
        """

        self.logger.info(
            f"Connecting to NATS {self.config.nats_server} with runner name {self.config.runner_name}..."
        )

        # connect to NATS server and jetstream
        await self.nc.connect(self.config.nats_server, name=self.config.runner_name)
        self.js = self.nc.jetstream()

        for workflow, _ in self.subjects.items():
            stream = (
                f"{self.config.krt_runtime_id}-{self.config.krt_version}-{workflow}"
            )
            subjects = [f"{stream}.*"]
            input_subject = f"{stream}.{self.config.runner_name}"

            # Create NATS stream
            await self.js.add_stream(name=stream, subjects=subjects)
            self.logger.info(f"Created stream {stream} with subjects: {subjects}")

            # Create NATS subscription for the entrypoint
            self.streams[workflow] = stream
            self.subscriptions[workflow] = await self.js.subscribe(
                stream=stream,
                subject=input_subject,
                config=ConsumerConfig(
                    deliver_policy=DeliverPolicy.NEW,
                ),
            )
            self.logger.info(
                f"Workflow {workflow} subscribed to NATS subject: '{input_subject}' from stream: '{stream}'"
            )

            # for each subscription start listening for messages
            await self.start_listening_messages(self.subscriptions[workflow], workflow)

    async def wait_for_response(self, request_id: str) -> None:
        while request_id in self.grpc_streams:
            await asyncio.sleep(0.1)
            pass

    async def start_listening_messages(self, sub, workflow) -> None:
        while True:
            # wait for the response
            self.logger.info(f"Waiting for response message...")
            msg = await sub.next_msg(timeout=self.config.request_timeout)
            await msg.ack()
            
            # prepare the grpc response message
            kre_nats_message = self.create_grpc_response(workflow, msg.data)
            response = self.make_response_object(workflow, kre_nats_message)

            self.logger.info(f"Response message received: {kre_nats_message.payload} with request_id {kre_nats_message.reply}")

            await self.response_to_grcp_stream(response, kre_nats_message.reply)
            # await grpc_stream.send_message(response)

    def create_kre_request_message(
        self, raw_msg: bytes, start: str, request_id: str
    ) -> bytes:

        """
        Creates a KreNatsMessage that packages the grpc request (raw_msg) and adds the required
        info needed to send the request to the NATS server.
        It returns the message in bytes, so it can be directly sent to the NATS server.

        :param raw_msg: the raw grpc request message
        :param start: the start time of the request

        :return: the message in bytes
        """

        tracking_id = str(uuid.uuid4())

        request_msg = KreNatsMessage()
        request_msg.tracking_id = tracking_id
        request_msg.payload.Pack(raw_msg)
        request_msg.reply = request_id
        t = request_msg.tracking.add()
        t.node_name = self.config.runner_name
        t.start = start
        t.end = datetime.utcnow().isoformat()

        return self._prepare_nats_request(request_msg.SerializeToString())

    def create_grpc_response(
        self, workflow: str, message_data: bytes
    ) -> KreNatsMessage:

        """
        Creates a gRPC response from the message data received from the NATS server.

        :param workflow: the workflow name
        :param message_data: the message data received from the NATS server

        :return: the gRPC response message
        """

        response_data = self._prepare_nats_response(message_data)

        response_msg = KreNatsMessage()
        response_msg.ParseFromString(response_data)

        if response_msg.error != "":
            self.logger.error(f"Received error message: {response_msg.error}")

            raise GRPCError(Status.INTERNAL, response_msg.error)

        return response_msg

    async def response_to_grcp_stream(self, response: bytes, request_id: str) -> None:
        """
        Sends the response to the gRPC stream.
        :param response: the response to be sent to the gRPC stream.
        :param stream: the gRPC stream.
        """
        grpc_stream = self.grpc_streams.pop(request_id)
        self.logger.info(f"{grpc_stream.peer}")

        await grpc_stream.send_message(response)
        # try:
        #     await grpc_stream.send_message(response)
        # except Exception as e:
        #     self.logger.error(f"Error sending response to gRPC stream: {e}")
            
        self.logger.info(f"gRPC response successfully sent")

    async def process_grpc_message(self, grpc_stream: Stream, workflow: str, request_id: str) -> None:

        """
        This function is called each time a gRPC message is received.

        It processes the message by sending it by the NATS server and waits for a response.

        :param grpc_stream: the gRPC stream
        :param workflow: the workflow name
        """

        start = datetime.utcnow().isoformat()

        try:
            # as multiple requests can be sent to the same workflow, we need to track each open gRPC stream to
            # send the response to the correct gRPC stream
            self.grpc_streams[request_id] = grpc_stream
            
            self.logger.info("############################# RECV_MESSAGE #############################")
            grpc_raw_msg = await grpc_stream.recv_message()
            self.logger.info(f"gRPC message received {grpc_raw_msg} from {grpc_stream.peer} and request_id {request_id}")

            # get the correct subject, subscription and stream depending on the workflow
            subject = self.subjects[workflow]
            subscription = self.subscriptions[workflow]
            stream = self.streams[workflow]

            # creates the msg to be sent to the NATS server
            request_msg = self.create_kre_request_message(
                grpc_raw_msg, start, request_id
            )

            # publish the msg to the NATS server
            await self.js.publish(stream=stream, subject=subject, payload=request_msg)
            self.logger.info(f"Message published to NATS subject: '{subject}' from stream: '{stream}'")
    
            # wait for the response
            #self.logger.info(f"Waiting for response message...")
            #msg = await subscription.next_msg(timeout=self.config.request_timeout)
            #await msg.ack()
            #self.logger.info(f"Response message received: {msg.data}")

            ## prepare the grpc response message
            #kre_nats_message = self.create_grpc_response(workflow, msg.data)
            #response = self.make_response_object(workflow, kre_nats_message)

            #self.logger.info(f"Response to request_id {kre_nats_message.reply}: {kre_nats_message.payload}")
            #await self.response_to_grcp_stream(response, kre_nats_message.reply)
            ## await grpc_stream.send_message(response)

        except Exception as err:
            err_msg = f"Exception on gRPC call : {err}"
            self.logger.error(err_msg)
            traceback.print_exc()

            if isinstance(err, GRPCError):
                raise err

    def _prepare_nats_request(self, msg: bytes) -> bytes:

        """
        Prepares the message to be sent to the NATS server by compressing it if needed.

        :param msg: the message to be sent to the NATS server.

        :return: the message to be sent to the NATS server compressed.
        """

        if len(msg) <= MESSAGE_THRESHOLD:
            return msg

        out = gzip.compress(msg, compresslevel=COMPRESS_LEVEL)

        if len(out) > MESSAGE_THRESHOLD:
            raise Exception("compressed message exceeds maximum size allowed of 1 MB.")

        self.logger.info(
            f"Original message size: {size_in_kb(msg)}. Compressed: {size_in_kb(out)}"
        )

        return out

    @staticmethod
    def _prepare_nats_response(msg: bytes) -> bytes:
        if msg.startswith(GZIP_HEADER):
            return gzip.decompress(msg)

        return msg


def size_in_kb(s: bytes) -> str:
    return f"{(len(s) / 1024):.2f} KB"
