from nats.js import JetStreamContext

from config import Config


class ContextObjectStore:
    def __init__(self, config: Config, js: JetStreamContext, logger):
        self.__config__ = config
        self.__logger__ = logger
        self.__object_store__ = js.object_store(config.nats_object_store)

    async def store_object(self, obj_name: str, data: bytes):
        pass

    async def get_object(self, obj_name: str) -> any:
        pass

    async def delete_object(self, obj_name: str):
        pass

