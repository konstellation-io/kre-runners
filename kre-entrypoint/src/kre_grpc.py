import abc
from datetime import datetime
import gzip
import traceback
import uuid

from grpclib import GRPCError
from grpclib.const import Status

from kre_nats_msg_pb2 import KreNatsMessage
from kre_measurements import KreMeasurements

COMPRESS_LEVEL = 9
MESSAGE_THRESHOLD = 1024 * 1024
GZIP_HEADER = b"\x1f\x8b"


# NOTE: EntrypointKRE will be extended by Entrypoint class auto-generated
class EntrypointKRE:
    def __init__(self, logger, nc, subjects, config):
        self.logger = logger
        self.subjects = subjects
        self.nc = nc
        self.config = config

    @abc.abstractmethod
    async def make_response_object(self, subject, response):
        # To be implemented on autogenerated entrypoint
        pass

    async def process_message(self, stream, subject) -> None:
        start = datetime.utcnow().isoformat()
        tracking_id = str(uuid.uuid4())

        try:
            raw_msg = await stream.recv_message()

            self.logger.info(f"gRPC message received")
            request_msg = KreNatsMessage()
            request_msg.tracking_id = tracking_id
            request_msg.payload.Pack(raw_msg)
            t = request_msg.tracking.add()
            t.node_name = "entrypoint"
            t.start = start
            t.end = datetime.utcnow().isoformat()

            nats_subject = self.subjects[subject]
            self.logger.info(f"Starting request/reply on NATS subject: '{nats_subject}'")

            nats_message = self._prepare_nats_request(request_msg.SerializeToString())

            nats_reply = await self.nc.request(
                nats_subject, nats_message, timeout=self.config.request_timeout
            )

            response_data = self._prepare_nats_response(nats_reply.data)

            self.logger.info(f"creating a response from message reply")
            response_msg = KreNatsMessage()
            response_msg.ParseFromString(response_data)

            if response_msg.error != "":
                self.logger.error(f"received error message: {response_msg.error}")

                raise GRPCError(Status.INTERNAL, response_msg.error)

            response = self.make_response_object(subject, response_msg)

            await stream.send_message(response)

            self.logger.info(f"gRPC successfully response")

        except Exception as err:
            err_msg = f"Exception on gRPC call : {err}"
            self.logger.error(err_msg)
            traceback.print_exc()

            if isinstance(err, GRPCError):
                raise err

    def _prepare_nats_request(self, msg: bytes) -> bytes:
        if len(msg) <= MESSAGE_THRESHOLD:
            return msg

        out = gzip.compress(msg, compresslevel=COMPRESS_LEVEL)

        if len(out) > MESSAGE_THRESHOLD:
            raise Exception("compressed message exceeds maximum size allowed of 1 MB.")

        self.logger.info(f"Original message size: {size_in_kb(msg)}. Compressed: {size_in_kb(out)}")

        return out

    def _prepare_nats_response(self, msg: bytes) -> bytes:
        if msg.startswith(GZIP_HEADER):
            return gzip.decompress(msg)

        return msg


def size_in_kb(s: bytes) -> str:
    return f"{(len(s) / 1024):.2f} KB"
