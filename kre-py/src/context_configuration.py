from __future__ import annotations

from logging import Logger

from nats.js.client import JetStreamContext
from nats.js.kv import KeyValue

from config import Config

ScopeProject = "project"
ScopeWorkflow = "workflow"
ScopeNode = "node"
ScopeUndefined = "undefined"


async def new_context_configuration(
    config: Config, logger: Logger, js: JetStreamContext
) -> ContextConfiguration:

    kv_stores_map: dict[str, KeyValue] = {}

    try:
        kv_store = await js.key_value(config.nats_key_value_store_project)
        kv_stores_map[ScopeProject] = kv_store

        kv_store = await js.key_value(config.nats_key_value_store_workflow)
        kv_stores_map[ScopeWorkflow] = kv_store

        kv_store = await js.key_value(config.nats_key_value_store_node)
        kv_stores_map[ScopeNode] = kv_store

    except Exception as err:
        logger.error(f"Error while getting the key value store: {err}")
        raise err

    return ContextConfiguration(config, logger, kv_stores_map)


class ContextConfiguration:
    """
    Provides a way to manipulate the jetstream key value store.
    """

    def __init__(
        self,
        config: Config,
        logger: Logger,
        kv_stores_map: dict[str, KeyValue],
    ):
        self.__config__: Config = config
        self.__logger__: Logger = logger
        self.__kv_stores_map__: dict[str, KeyValue] = kv_stores_map

    async def set(self, key: str, value: str, scope: str = ScopeNode) -> None:
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

    async def get(self, key: str, scope: str = ScopeUndefined) -> str:
        """
        Gets a value from the key value store by given scope.
        If no scope is given, the default scope (node) is used.
        If no value is found using the default scope, the search continues upwards.
        """

        # search by scope
        if scope is not ScopeUndefined:
            try:
                kv_store = self.__kv_stores_map__[scope]
                entry = await kv_store.get(key)
                if entry.value is not None:
                    return entry.value.decode("utf-8")
                else:
                    raise Exception(f"No value found")

            except Exception as err:
                self.__logger__.error(f"Error while getting the value for key {key}: {err}")
                raise err

        # default search
        all_scopes_in_order = [ScopeNode, ScopeWorkflow, ScopeProject]
        for scope in all_scopes_in_order:
            try:
                kv_store = self.__kv_stores_map__[scope]
                entry = await kv_store.get(key)
                if entry.value is not None:
                    return entry.value.decode("utf-8")

            except Exception as err:
                self.__logger__.error(f"Error while getting the value for key {key}: {err}")
                continue

        self.__logger__.error(
            f"Error while getting the value for key {key} in all scopes: No value found"
        )
        raise Exception("No value found")

    async def delete(self, key: str, scope: str = ScopeNode) -> None:
        """
        Deletes a value from the key value store by given scope.
        If no scope is given, the default scope (node) is used.
        """

        try:
            kv_store = self.__kv_stores_map__[scope]
            await kv_store.delete(key)

        except Exception as err:
            self.__logger__.error(f"Error while deleting the value for key {key}: {err}")
            raise err
