import os
from logging import Logger
from typing import Any, Awaitable, Callable

from google.protobuf.message import Message as ProtobufMessage
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from pymongo import MongoClient

from config import Config
from context_configuration import ContextConfiguration
from context_data import ContextData
from context_measurement import ContextMeasurement
from context_objectstore import ContextObjectStore
from context_prediction import ContextPrediction
from kre_nats_msg_pb2 import EARLY_EXIT, EARLY_REPLY, ERROR, OK, KreNatsMessage, MessageType

PublishMsgFunc = Callable[[ProtobufMessage, Any, MessageType, str], Awaitable]
PublishAnyFunc = Callable[[ProtobufMessage, Any, MessageType, str], Awaitable]

DEFAULT_CHANNEL = ""


class HandlerContext:
    def __init__(
        self,
        config: Config,
        nc: NatsClient,
        js: JetStreamContext,
        mongo_conn: MongoClient,
        logger: Logger,
        publish_msg: PublishMsgFunc,
        publish_any: PublishAnyFunc,
        configuration: ContextConfiguration,
    ):
        self.__config__: Config = config
        self.__publish_msg__ = publish_msg
        self.__publish_any__ = publish_any
        self.__request_msg__: KreNatsMessage
        self.logger = logger
        self.prediction = ContextPrediction(config, nc, logger)
        self.measurement = ContextMeasurement(config, logger)
        self.db = ContextData(config, nc, js, mongo_conn, logger)
        self.configuration = configuration

        if config.nats_object_store is not None:
            self.object_store = ContextObjectStore(config, logger)

    def path(self, relative_path):
        return os.path.join(self.__config__.base_path, relative_path)

    def set_request_msg(self, request_msg):
        self.__request_msg__ = request_msg

    def get_request_id(self) -> str:
        """get_request_id will return the payload's original request ID."""
        return str(self.__request_msg__.request_id)

    async def send_output(self, message: ProtobufMessage, channel: str = DEFAULT_CHANNEL) -> None:
        """
        send_output will send a desired typed proto payload to the node's subject.
        By specifying a channel, the message will be sent to that subject's subtopic.

        send_output converts the proto message into an any type. This means the following node will recieve
        an any type protobuf.

        GRPC requests can only be answered once. So once the entrypoint has been replied by the exitpoint,
        all following replies to the entrypoint from the same request will be ignored.

        :param  message:  The message to publish.
        :param channel (optional):  The subsubject channel where the message will be published.
        """
        await self.__publish_msg__(message, self.__request_msg__, OK, channel)

    async def send_any(self, message: ProtobufMessage, channel: str = DEFAULT_CHANNEL) -> None:
        """
        send_any will send any type of proto payload to the node's subject.
        By specifying a channel, the message will be sent to that subject's subtopic.

        As a difference from SendOutput, send_any will not convert your proto structure.

        Use this function when you wish to simply redirect your node's payload without unpackaging.
        Once the entrypoint has been replied, all following replies to the entrypoint will be ignored.

        :param  message:  The message to publish.
        :param channel (optional):  The subsubject channel where the message will be published.
        """
        await self.__publish_any__(message, self.__request_msg__, OK, channel)

    async def send_early_reply(
        self, message: ProtobufMessage, channel: str = DEFAULT_CHANNEL
    ) -> None:
        """
        send_early_reply works as the SendOutput functionality
        with the addition of typing this message as an early reply.

        :param  message:  The message to publish.
        :param channel (optional):  The subsubject channel where the message will be published.
        """
        await self.__publish_msg__(message, self.__request_msg__, EARLY_REPLY, channel)

    async def send_early_exit(
        self, message: ProtobufMessage, channel: str = DEFAULT_CHANNEL
    ) -> None:
        """
        send_early_exit works as the SendOutput functionality
        with the addition of typing this message as an early exit.

        :param  message:  The message to publish.
        :param channel (optional):  The subsubject channel where the message will be published.
        """
        await self.__publish_msg__(message, self.__request_msg__, EARLY_EXIT, channel)

    def is_message_ok(self) -> bool:
        """is_message_ok returns true if incoming message is of message type OK."""
        return self.__request_msg__.message_type == OK

    def is_message_error(self) -> bool:
        """is_message_error returns true if incoming message is of message type ERROR."""
        return self.__request_msg__.message_type == ERROR

    def is_message_early_reply(self) -> bool:
        """is_message_early_reply returns true if incoming message is of message type EARLY_REPLY."""
        return self.__request_msg__.message_type == EARLY_REPLY

    def is_message_early_exit(self) -> bool:
        """is_message_early_exit returns true if incoming message is of message type EARLY_EXIT."""
        return self.__request_msg__.message_type == EARLY_EXIT
