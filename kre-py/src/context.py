import os

from context_measurement import ContextMeasurement
from context_prediction import ContextPrediction
from context_data import ContextData


class HandlerContext:
    def __init__(self, config, nc, mongo_conn, logger, early_reply_func):
        self.__data__ = lambda: None
        self.__config__ = config
        self.__early_reply__ = early_reply_func
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
        if self.__request_msg__.replied:
            raise Exception("error the message was replied previously")

        self.__request_msg__.replied = True
        await self.__early_reply__(response, self.__request_msg__.request_id)

    def early_exit(self):
        self.__request_msg__.early_exit = True
