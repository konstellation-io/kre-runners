import sys
from logging import Logger

from nats.js import JetStreamContext
from nats.js.api import ObjectMeta
from nats.js.object_store import ObjectStore

from config import Config


class ContextObjectStore:
    """
    Provides a way to manipulate objects stored in Object Store.
    """

    def __init__(self, config: Config, js: JetStreamContext, logger: Logger):
        self.__config__: Config = config
        self.__logger__: Logger = logger
        self.__object_store__: ObjectStore | None = None

        self.__connect(js)

    def __connect(self, js: JetStreamContext):
        """
        Connects to the Object Store.

        :param js: JetStream client.
        :raises Exception: If there is an error while connecting to Object Store.
        """

        self.__logger__.info(f"Connecting to Object Store {self.__config__.nats_object_store}...")
        try:
            self.__object_store__: ObjectStore = await js.object_store(self.__config__.nats_object_store)
        except Exception as e:
            self.__logger__.error(f"Error binding the object store: {e}")
            sys.exit(1)

    async def store_object(self, key: str, payload: bytes):
        """
        Stores a payload with the desired key to Object Store.

        :param key: the object name.
        :param payload: a sequence of bytes.
        :raises Exception: If the Object Store does not exist.
        :raises Exception: If the payload is empty or null.
        :raises Exception: If there is an error while storing the object.
        """

        if self.__object_store__ is None:
            raise Exception("the object store does not exist")

        if not payload:
            raise Exception("the payload cannot be empty")

        obj_meta = ObjectMeta(name=key)
        try:
            await self.__object_store__.put(meta=obj_meta, data=payload)
        except Exception as e:
            raise Exception(f"error storing object with key {key} to the object store: {e}")

        self.__logger__.debug(
            f"File with key {key} successfully stored in object store {self.__config__.nats_object_store}"
        )

    async def get_object(self, key: str) -> bytes:
        """
        Retrieves a payload with the desired key from Object Store.

        :param key: the object name.
        :returns: a sequence of bytes.
        :raises Exception: If the Object Store does not exist.
        :raises Exception: If there is an error while retrieving the object.
        """

        if self.__object_store__ is None:
            raise Exception("the object store does not exist")

        try:
            obj_result: ObjectStore.ObjectResult = await self.__object_store__.get(name=key)
        except Exception as e:
            raise Exception(f"error retrieving object with key {key} from the object store: {e}")

        self.__logger__.debug(
            f"File with key {key} successfully retrieved from object store {self.__config__.nats_object_store}"
        )
        return obj_result.data

    async def delete_object(self, key: str):
        """
        Deletes an object from Object Store.

        :param key: the object name.
        :raises Exception: If the Object Store does not exist.
        :raises Exception: If there is an error while deleting the object.
        """
        raise NotImplementedError

