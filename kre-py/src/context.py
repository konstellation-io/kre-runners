import os
from typing import Callable, Any, Awaitable

from google.protobuf.message import Message

from context_measurement import ContextMeasurement
from context_prediction import ContextPrediction
from context_data import ContextData
from kre_nats_msg_pb2 import OK, MessageType, EARLY_REPLY, EARLY_EXIT

PublishMsgFunc = Callable[[Message, Any, MessageType, str], Awaitable]
PublishAnyFunc = Callable[[Message, Any, MessageType, str], Awaitable]

DEFAULT_CHANNEL = ""


class HandlerContext:
    def __init__(self, config, nc, mongo_conn, logger, publish_msg: PublishMsgFunc, publish_any: PublishAnyFunc):
        self.__data__ = lambda: None
        self.__config__ = config
        self.__publish_msg__ = publish_msg
        self.__publish_any__ = publish_any
        self.__request_msg__ = None
        self.logger = logger
        self.prediction = ContextPrediction(config, nc, logger)
        self.measurement = ContextMeasurement(config, logger)
        self.db = ContextData(config, nc, mongo_conn, logger)

    def path(self, relative_path):
        return os.path.join(self.__config__.base_path, relative_path)

    def set(self, key: str, value: any):
        setattr(self.__data__, key, value)

    def get(self, key: str) -> any:
        return getattr(self.__data__, key)

    async def early_reply(self, response):
        await self.__publish_msg__(response, self.__request_msg__, EARLY_REPLY, DEFAULT_CHANNEL)

    def early_exit(self, message):
        self.__publish_msg__(message, self.__request_msg__, EARLY_EXIT, DEFAULT_CHANNEL)

    def set_request_msg(self, request_msg):
        self.__request_msg__ = request_msg

    async def send_output(self, message: Message, channel: str = DEFAULT_CHANNEL):
        await self.__publish_msg__(message, self.__request_msg__, OK, channel)

    async def send_any(self, message, channel: str = DEFAULT_CHANNEL):
        await self.__publish_any__(message, self.__request_msg__, OK, channel)

    def is_request_early_reply(self):
        return self.__request_msg__.message_type == EARLY_REPLY
