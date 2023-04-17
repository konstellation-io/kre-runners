import abc
import gzip
import traceback

from exceptions.exceptions import CompressedMessageTooLargeException

from grpclib.server import Stream
from grpclib import GRPCError
from grpclib.const import Status
from nats.js.api import ConsumerConfig, DeliverPolicy
from nats.aio.client import Client as NATS

from kre_nats_msg_pb2 import KreNatsMessage, OK

COMPRESS_LEVEL = 9
GZIP_HEADER = b"\x1f\x8b"


# NOTE: EntrypointKRE will be extended by Entrypoint class auto-generated
class EntrypointKRE:
    def __init__(self, logger, jetstream_data, config):
        self.logger = logger
        self.nc = NATS()
        self.js = None
        self.jetstream_data = jetstream_data
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
            f"Connecting to NATS {self.config.nats_server} "
            f"with runner name {self.config.runner_name}..."
        )

        await self.nc.connect(self.config.nats_server, name=self.config.runner_name)
        self.js = self.nc.jetstream()

    async def process_grpc_message(
        self, grpc_stream: Stream, workflow: str, request_id: str
    ) -> None:
        """
        This function is called each time a gRPC message is received.

        It processes the message by sending it by the NATS server and waits for a response.

        :param grpc_stream: The gRPC stream.
        :param workflow: The workflow name.
        :param request_id: The gRPC request id.
        """

        try:
            # As multiple requests can be sent to the same workflow, we need to track each
            # open gRPC stream to send the response to the correct gRPC stream
            if "grpc_streams" not in self.jetstream_data[workflow]:
                self.jetstream_data[workflow]["grpc_streams"] = {request_id: grpc_stream}
            else:
                self.jetstream_data[workflow]["grpc_streams"][request_id] = grpc_stream

            grpc_raw_msg = await grpc_stream.recv_message()
            self.logger.debug(
                f"gRPC message received {grpc_raw_msg} "
                f"from {grpc_stream.peer} and request_id {request_id}"
            )

            # Get the correct subject, subscription and stream depending on the workflow
            stream = self.jetstream_data[workflow]["stream"]
            input_subject = self.jetstream_data[workflow]["input_subject"]
            output_subject = self.jetstream_data[workflow]["output_subject"]

            stream_info = await self.js.stream_info(stream)
            stream_max_size = stream_info.config.max_msg_size or -1
            server_max_size = self.nc.max_payload

            max_msg_size = (
                min(stream_max_size, server_max_size) if stream_max_size != -1 else server_max_size
            )

            # Creates the msg to be sent to the NATS server
            request_msg = self._create_kre_request_message(grpc_raw_msg, request_id, max_msg_size)

            # Create an ephemeral subscription for the request
            sub = await self.js.subscribe(
                stream=stream,
                subject=input_subject,
                config=ConsumerConfig(
                    deliver_policy=DeliverPolicy.NEW,
                    ack_wait=22 * 3600,  # 22 hours
                ),
                manual_ack=True,
            )

            await self.js.publish(stream=stream, subject=output_subject, payload=request_msg)
            self.logger.info(
                f"Message published to NATS subject: '{output_subject}' from stream: '{stream}'"
            )

            # Wait until a message for the request arrives ignoring the rest
            message_recv = False

            while not message_recv:
                self.logger.info("Waiting for response message...")
                msg = await sub.next_msg(timeout=self.config.request_timeout)
                self.logger.debug(f"Response message received: {msg.data}")

                kre_nats_message = self._create_grpc_response(msg.data)

                if kre_nats_message.request_id == request_id:
                    message_recv = True
                    await sub.unsubscribe()
                    response = self.make_response_object(workflow, kre_nats_message)
                    await self._respond_to_grpc_stream(
                        response, workflow, kre_nats_message.request_id
                    )

                await msg.ack()

        except Exception as err:
            err_msg = f"Exception on gRPC call : {err}"
            self.logger.error(err_msg)
            traceback.print_exc()

            if isinstance(err, GRPCError):
                raise err

    def _create_kre_request_message(self, raw_msg: bytes, request_id: str, max_msg_size: int) -> bytes:
        """
        Creates a KreNatsMessage that packages the grpc request (raw_msg) and adds the required
        info needed to send the request to the NATS server.
        It returns the message in bytes, so it can be directly sent to the NATS server.

        :param raw_msg: the raw grpc request message
        :param request_id: the id of the gRPC that should be responded.

        :return: the message in bytes
        """
        request_msg = KreNatsMessage()
        request_msg.payload.Pack(raw_msg)
        request_msg.request_id = request_id
        request_msg.from_node = self.config.krt_node_name
        request_msg.message_type = OK
        return self._prepare_nats_request(request_msg.SerializeToString(), max_msg_size)

    def _create_grpc_response(self, message_data: bytes) -> KreNatsMessage:
        """
        Creates a gRPC response from the message data received from the NATS server.

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

    async def _respond_to_grpc_stream(
        self, response: bytes, workflow: str, request_id: str
    ) -> None:
        """
        Sends a response to the corresponding gRPC stream based on the request id.

        :param response: the response to be sent to the gRPC stream.
        :param workflow: the workflow name
        :param request_id: the gRPC request id that should be responded.
        """

        grpc_stream = self.jetstream_data[workflow]["grpc_streams"].pop(request_id)
        await grpc_stream.send_message(response)
        self.logger.info(f"gRPC request '{request_id}' response successfully sent")

    def _prepare_nats_request(self, msg: bytes, max_msg_size: int) -> bytes:
        """
        Prepares the message to be sent to the NATS server by compressing it if needed.

        :param msg: the message to be sent to the NATS server.

        :return: the message to be sent to the NATS server compressed.
        """
        if len(msg) <= max_msg_size:
            return msg

        out = gzip.compress(msg, compresslevel=COMPRESS_LEVEL)

        if len(out) > max_msg_size:
            raise CompressedMessageTooLargeException("compressed message exceeds maximum size allowed.")

        self.logger.debug(
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
