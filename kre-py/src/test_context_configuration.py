from logging import Logger
from unittest.mock import AsyncMock, MagicMock

import pytest
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from nats.js.kv import KeyValue

from config import Config
from context_configuration import ContextConfiguration, Scope, new_context_configuration


@pytest.fixture()
def nats_mock() -> MagicMock:
    nats_mock = MagicMock(NatsClient)

    return nats_mock


@pytest.fixture()
def jetstream_mock() -> MagicMock:
    jetstream_mock = MagicMock(JetStreamContext)

    return jetstream_mock


@pytest.fixture()
def config_mock() -> MagicMock:
    config = MagicMock(Config)

    return config


@pytest.fixture()
def logger_mock() -> MagicMock:
    logger = MagicMock(Logger)

    return logger


@pytest.fixture()
def kv_map_mock() -> MagicMock:
    kv_map = MagicMock(dict[Scope, KeyValue])

    return kv_map


@pytest.fixture()
def context_configuration(
    config_mock: MagicMock, logger_mock: MagicMock, kv_map_mock: MagicMock
) -> ContextConfiguration:
    context_configuration = ContextConfiguration(config_mock, logger_mock, kv_map_mock)
    return context_configuration
