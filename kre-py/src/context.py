import os
from typing import Callable, Any, Awaitable

from google.protobuf.message import Message

from context_measurement import ContextMeasurement
from context_prediction import ContextPrediction
from context_data import ContextData
from context_objectstore import ContextObjectStore
from kre_nats_msg_pb2 import ERROR, OK, MessageType, EARLY_REPLY, EARLY_EXIT

PublishMsgFunc = Callable[[Message, Any, MessageType, str], Awaitable]
PublishAnyFunc = Callable[[Message, Any, MessageType, str], Awaitable]

DEFAULT_CHANNEL = ""


class HandlerContext:
    def __init__(self, config, nc, js, mongo_conn, logger, publish_msg: PublishMsgFunc, publish_any: PublishAnyFunc):
        self.__data__ = lambda: None
        self.__config__ = config
        self.__publish_msg__ = publish_msg
        self.__publish_any__ = publish_any
        self.__request_msg__ = None
        self.logger = logger
        self.prediction = ContextPrediction(config, nc, logger)
        self.measurement = ContextMeasurement(config, logger)
        self.db = ContextData(config, nc, mongo_conn, logger)

        if config.nats_object_store is not None:
            self.object_store = ContextObjectStore(config, js, logger)

    def path(self, relative_path):
        return os.path.join(self.__config__.base_path, relative_path)

    def set(self, key: str, value: any):
        setattr(self.__data__, key, value)

    def get(self, key: str) -> any:
        return getattr(self.__data__, key)

    def set_request_msg(self, request_msg):
        self.__request_msg__ = request_msg

    def get_request_id(self) -> str:
        """get_request_id will return the payload's original request ID."""
        return self.__request_msg__.request_id

    async def send_output(self, message: Message, channel: str = DEFAULT_CHANNEL):
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

    async def send_any(self, message, channel: str = DEFAULT_CHANNEL):
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

    async def send_early_reply(self, message, channel: str = DEFAULT_CHANNEL):
        """
        send_early_reply works as the SendOutput functionality
        with the addition of typing this message as an early reply.

        :param  message:  The message to publish.
        :param channel (optional):  The subsubject channel where the message will be published.
        """
        await self.__publish_msg__(message, self.__request_msg__, EARLY_REPLY, channel)

    async def send_early_exit(self, message, channel: str = DEFAULT_CHANNEL):
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
