import os

from context_measurement import ContextMeasurement
from context_prediction import ContextPrediction
from context_data import ContextData


class HandlerContext:
    def __init__(self, config, nc, mongo_conn, logger, reply, send_output):
        self.__data__ = lambda: None
        self.__config__ = config
        self.__reply__ = reply
        self.__send_output__ = send_output
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

    async def reply(self, response):
        if self.__request_msg__.replied:
            raise Exception("error the message was replied previously")

        self.__request_msg__.replied = True
        await self.__reply__(self.__request_msg__.reply, response)

    async def send_output(self, response):
        await self.__send_output__(self.__config__.nats_output, self.__config__.nats_input, response)
