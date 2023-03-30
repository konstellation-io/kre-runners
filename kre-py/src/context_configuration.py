from __future__ import annotations

from enum import Enum
from logging import Logger

from nats.js.client import JetStreamContext
from nats.js.kv import KeyValue

from config import Config


class Scope(Enum):
    """
    Defines the scope of the key value store.
    """

    PROJECT = "project"
    WORKFLOW = "workflow"
    NODE = "node"
    UNDEFINED = "undefined"


class ContextConfiguration:
    """
    Provides a way to manipulate the jetstream key value store.
    """

    async def new_context_configuration(
        self, config: Config, logger: Logger, js: JetStreamContext
    ) -> ContextConfiguration:
        self.__config__: Config = config
        self.__logger__: Logger = logger
        self.__js__: JetStreamContext = js
        self.__kv_stores_map__: dict[Scope, KeyValue] = {}

        await self.init_kv_stores_map()
        return self

    async def init_kv_stores_map(self) -> None:
        """
        Initializes the kv_stores_map attribute.
        """

        try:
            kv_store = await self.__js__.key_value(self.__config__.nats_key_value_store_project)
            self.__kv_stores_map__[Scope.PROJECT] = kv_store

            kv_store = await self.__js__.key_value(self.__config__.nats_key_value_store_workflow)
            self.__kv_stores_map__[Scope.WORKFLOW] = kv_store

            kv_store = await self.__js__.key_value(self.__config__.nats_key_value_store_node)
            self.__kv_stores_map__[Scope.NODE] = kv_store

        except Exception as err:
            self.__logger__.error(f"Error while getting the key value store: {err}")
            raise err

    async def set(self, key: str, value: str, scope: Scope = Scope.NODE) -> None:
        """
        Sets a value in the key value store by given scope.
        If no scope is given, the default scope (node) is used.
        """

        try:
            kv_store = self.__kv_stores_map__[scope]
            await kv_store.put(key, bytes(value, "utf-8"))

        except Exception as err:
            self.__logger__.error(f"Error while setting the value for key {key}: {err}")
            raise err

    async def get(self, key: str, scope: Scope = Scope.UNDEFINED) -> str:
        """
        Gets a value from the key value store by given scope.
        If no scope is given, the default scope (node) is used.
        If no value is found using the default scope, the search continues upwards.
        """

        # search by scope
        if scope is not Scope.UNDEFINED:
            try:
                kv_store = self.__kv_stores_map__[scope]
                entry = await kv_store.get(key)
                if entry.value is not None:
                    return str(entry.value)
                else:
                    raise Exception(f"Error getting the value for key {key}: no value found")

            except Exception as err:
                self.__logger__.error(f"Error while getting the value for key {key}: {err}")
                raise err

        # default search
        all_scopes_in_order = [Scope.NODE, Scope.WORKFLOW, Scope.PROJECT]
        for scope in all_scopes_in_order:
            try:
                kv_store = self.__kv_stores_map__[scope]
                entry = await kv_store.get(key)
                if entry.value is not None:
                    return str(entry.value)

            except Exception as err:
                self.__logger__.error(f"Error while getting the value for key {key}: {err}")
                continue

        raise Exception(f"Error while getting the value for key {key}: no value found")

    def delete(self, key: str, scope: Scope = Scope.NODE) -> None:
        """
        Deletes a value from the key value store by given scope.
        If no scope is given, the default scope (node) is used.
        """

        try:
            kv_store = self.__kv_stores_map__[scope]
            kv_store.delete(key)

        except Exception as err:
            self.__logger__.error(f"Error while deleting the value for key {key}: {err}")
            raise err
