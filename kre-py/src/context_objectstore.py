import asyncio
import shlex
import shutil
import subprocess
import sys
import tempfile
from logging import Logger

from config import Config


class ContextObjectStoreException(Exception):
    pass


class ContextObjectStore:
    """
    Provides a way to manipulate objects stored in Object Store.
    """
    __nats_bin__: str = "nats"

    def __init__(self, config: Config, logger: Logger):
        self.__config__: Config = config
        self.__logger__: Logger = logger

        base_cmd = f"{self.__nats_bin__} --server={config.nats_server} object"
        obj_store = config.nats_object_store

        # Command templates
        self.__info_cmd__ = f"{base_cmd} info {obj_store}"
        self.__put_obj_cmd__ = f"{base_cmd} put {obj_store} --name={{obj_name}} --force --no-progress"
        self.__get_obj_cmd__ = f"{base_cmd} get {obj_store} {{obj_name}} --output={{dst_path}} --force --no-progress"
        self.__del_obj_cmd__ = f"{base_cmd} del {obj_store} {{obj_name}} --force"

        self.__check_prereqs()

    def __check_prereqs(self):
        """
        Ensures some pre-requisites are met before the Object Store context can be used.

        :raises Exception: If the nats-cli binary cannot be found in the system PATH.
        :raises Exception: If the Object Store configured for this runner does not exist in JetStream.
        """

        self.__logger__.debug("Looking for the nats-cli binary...")
        if shutil.which(self.__nats_bin__) is None:
            self.__logger__.debug("Could not find nats-cli in the system PATH. Is it installed?")
            sys.exit(1)

        self.__logger__.info(f"Checking if the Object Store {self.__config__.nats_object_store} exists...")
        out = subprocess.run(args=self.__info_cmd__.split(), capture_output=True)
        if out.returncode != 0:
            self.__logger__.error(
                f"Error while getting info for Object Store {self.__config__.nats_object_store}: "
                f"{out.stderr.decode('utf-8')}"
            )
            sys.exit(1)

        self.__logger__.info(f"Successfully bound to Object Store {self.__config__.nats_object_store}")

    async def save(self, key: str, payload: bytes):
        """
        Stores a payload with the desired key to Object Store.

        :param key: the object name.
        :param payload: a sequence of bytes.
        :raises ContextObjectStoreException: If the payload is empty or null.
        :raises ContextObjectStoreException: If there is an error while storing the object.
        """

        if not payload:
            raise ContextObjectStoreException("the payload cannot be empty")

        cmd = self.__put_obj_cmd__.format(obj_name=shlex.quote(key))
        subp = await asyncio.create_subprocess_shell(
            cmd,
            stdin=asyncio.subprocess.PIPE,
            stdout=asyncio.subprocess.DEVNULL,
            stderr=asyncio.subprocess.PIPE,
        )
        _, stderr = await subp.communicate(input=payload)
        if subp.returncode != 0:
            raise ContextObjectStoreException(f"error storing object with key {key} to the object store: {stderr.decode('utf-8')}")

        self.__logger__.debug(
            f"File with key {key} successfully stored in object store {self.__config__.nats_object_store}"
        )

    async def get(self, key: str) -> bytes:
        """
        Retrieves a payload with the desired key from Object Store.

        :param key: the object name.
        :returns: a sequence of bytes
        :raises ContextObjectStoreException: If there is an error while retrieving the object.
        """

        with tempfile.NamedTemporaryFile(mode="rb") as fd:
            cmd = self.__get_obj_cmd__.format(obj_name=shlex.quote(key), dst_path=fd.name)
            subp = await asyncio.create_subprocess_shell(
                cmd,
                stdin=None,
                stdout=asyncio.subprocess.DEVNULL,
                stderr=asyncio.subprocess.PIPE,
            )
            _, stderr = await subp.communicate()
            if subp.returncode != 0:
                raise ContextObjectStoreException(
                    f"error retrieving object with key {key} from the object store: {stderr.decode('utf-8')}"
                )

            payload = fd.read()

        self.__logger__.debug(
            f"File with key {key} successfully retrieved from object store {self.__config__.nats_object_store}"
        )
        return payload

    async def delete(self, key: str):
        """
        Deletes an object from Object Store.

        :param key: the object name.
        :raises ContextObjectStoreException: If there is an error while deleting the object.
        """
        cmd = self.__del_obj_cmd__.format(obj_name=shlex.quote(key))
        subp = await asyncio.create_subprocess_shell(
            cmd,
            stdin=None,
            stdout=asyncio.subprocess.DEVNULL,
            stderr=asyncio.subprocess.PIPE,
        )
        _, stderr = await subp.communicate()
        if subp.returncode != 0:
            raise ContextObjectStoreException(f"error deleting object with key {key} from the object store: {stderr.decode('utf-8')}")

        self.__logger__.debug(
            f"File with key {key} successfully deleted from object store {self.__config__.nats_object_store}"
        )