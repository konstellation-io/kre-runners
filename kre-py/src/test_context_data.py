import json
from unittest import mock

import pytest

from datetime import datetime
from unittest.mock import MagicMock, AsyncMock

from nats.aio.client import Client
from nats.aio.msg import Msg
from nats.errors import TimeoutError
from pymongo import MongoClient
from pymongo.collection import Collection

from context_data import ContextData

DATETIME_FORMAT = '%Y-%m-%d %H:%M:%S.%f'


@pytest.fixture()
def nats_mock():
    nats_mock = MagicMock(Client)

    return nats_mock


@pytest.fixture()
def mongodb_mock():
    mongodb_mock = MagicMock(MongoClient)

    return mongodb_mock


@pytest.fixture
def context_data(config, nats_mock, mongodb_mock, simple_logger):
    context_measurement = ContextData(config, nats_mock, mongodb_mock, simple_logger)
    return context_measurement


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_save_successful_response_expect_ok(context_data, config):
    msg = MagicMock(Msg)
    msg.data = bytes(json.dumps({"success": True}), encoding="utf-8")
    context_data.__nc__.request = AsyncMock(return_value=msg)
    await context_data.save("test_collection", {
        "Time": datetime.utcnow().strftime(DATETIME_FORMAT),
        "Result": "Tested",
        "TicketID": "1234",
        "Asset": "A12345C",
    })


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_save_unsuccessful_response_expect_ok(context_data, config, simple_logger):
    msg = MagicMock(Msg)
    msg.data = bytes(json.dumps({"success": False}), encoding="utf-8")
    context_data.__nc__.request = AsyncMock(return_value=msg)
    with mock.patch.object(simple_logger, 'error') as mock_logger_error:
        await context_data.save("test_collection", {
            "Time": datetime.utcnow().strftime(DATETIME_FORMAT),
            "Result": "Tested",
            "TicketID": "1234",
            "Asset": "A12345C",
        })

        mock_logger_error.assert_called_with("Unexpected error saving data")


@pytest.mark.unittest
@pytest.mark.parametrize("collection_name", [1234, ""])
@pytest.mark.asyncio
async def test_save_wrong_collection_type_expect_exception(context_data, config, collection_name):
    with pytest.raises(Exception):
        await context_data.save(collection_name, {
            "Time": datetime.utcnow().strftime(DATETIME_FORMAT),
            "Result": "Tested",
            "TicketID": "1234",
            "Asset": "A12345C",
        })


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_save_expect_timeout_exception(context_data, config, simple_logger):
    context_data.__nc__.request = AsyncMock(side_effect=TimeoutError)

    with mock.patch.object(simple_logger, 'error') as mock_logger_error:
        await context_data.save("test", {
            "Time": datetime.utcnow().strftime(DATETIME_FORMAT),
            "Result": "Tested",
            "TicketID": "1234",
            "Asset": "A12345C",
        })

        mock_logger_error.assert_called_with("Error saving data: request timed out")


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_find_expect_ok(context_data, config):
    coll = "test_collection"
    collection_mock = MagicMock(Collection)
    collection_mock.find = MagicMock(return_value=["response1", "response2"])
    context_data.__mongo_conn__ = {f"{config.mongo_data_db_name}": {coll: collection_mock}}
    response = await context_data.find(coll, {
        "TicketID": "1234",
        "Asset": "A5678",
    })

    assert len(response) == 2
    assert response[0] == "response1"
    assert response[1] == "response2"


@pytest.mark.unittest
@pytest.mark.parametrize("collection_name, query", [(1234, {
    "TicketID": "1234",
    "Asset": "A5678",
}), ("collection_test", "some wrong query string")])
@pytest.mark.asyncio
async def test_find_wrong_type_expect_exception(context_data, config, collection_name, query):
    with pytest.raises(Exception):
        await context_data.find(collection_name, query)


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_find_expect_exception(context_data, config):
    with pytest.raises(Exception):
        coll = "test_collection"
        collection_mock = MagicMock(Collection)
        collection_mock.find = MagicMock(side_effect=Exception)
        context_data.__mongo_conn__ = {f"{config.mongo_data_db_name}": {coll: collection_mock}}
        await context_data.find(coll, {
            "TicketID": "1234",
            "Asset": "A5678",
        })
