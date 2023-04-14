import json
from datetime import datetime
from unittest import mock
from unittest.mock import AsyncMock, MagicMock

import pytest
from nats.aio.client import Client as NatsClient
from nats.aio.msg import Msg
from nats.errors import TimeoutError
from nats.js.client import JetStreamContext
from pymongo import MongoClient
from pymongo.collection import Collection

from context_data import ContextData

DATETIME_FORMAT = "%Y-%m-%d %H:%M:%S.%f"


@pytest.fixture()
def nats_mock():
    nats_mock = MagicMock(NatsClient)

    return nats_mock


@pytest.fixture()
def jetstream_mock():
    jetstream_mock = MagicMock(JetStreamContext)

    return jetstream_mock


@pytest.fixture()
def mongodb_mock():
    mongodb_mock = MagicMock(MongoClient)

    return mongodb_mock


@pytest.fixture
def context_data(config, nats_mock, jetstream_mock, mongodb_mock, simple_logger):
    context_measurement = ContextData(config, nats_mock, mongodb_mock, simple_logger)
    return context_measurement


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_save_successful_response_expect_ok(context_data, simple_logger, config):
    # GIVEN a successful response
    msg = MagicMock(Msg)
    msg.data = bytes(json.dumps({"success": True}), encoding="utf-8")
    context_data.__nc__.request = AsyncMock(return_value=msg)

    # WHEN the save method is called
    with mock.patch.object(simple_logger, "error") as mock_logger_error:
        await context_data.save(
            "test_collection",
            {
                "Time": datetime.utcnow().strftime(DATETIME_FORMAT),
                "Result": "Tested",
                "TicketID": "1234",
                "Asset": "A12345C",
            },
        )

        # THEN expect no error logs to be written
        mock_logger_error.assert_not_called()


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_save_unsuccessful_response_expect_ok(context_data, config, simple_logger):
    # GIVEN an unsuccessful response
    msg = MagicMock(Msg)
    msg.data = bytes(json.dumps({"success": False}), encoding="utf-8")
    context_data.__nc__.request = AsyncMock(return_value=msg)

    # WHEN the save method is called
    with mock.patch.object(simple_logger, "error") as mock_logger_error:
        await context_data.save(
            "test_collection",
            {
                "Time": datetime.utcnow().strftime(DATETIME_FORMAT),
                "Result": "Tested",
                "TicketID": "1234",
                "Asset": "A12345C",
            },
        )

        # THEN expect an error log is written with the given prompt
        mock_logger_error.assert_called_with("Unexpected error saving data")


@pytest.mark.unittest
@pytest.mark.parametrize("collection_name", [1234, ""])
@pytest.mark.asyncio
async def test_save_wrong_collection_type_expect_exception(context_data, config, collection_name):
    # GIVEN a wrong collection_name type
    # WHEN the save method is called,
    # THEN assert an exception is thrown
    with pytest.raises(Exception):
        await context_data.save(
            collection_name,
            {
                "Time": datetime.utcnow().strftime(DATETIME_FORMAT),
                "Result": "Tested",
                "TicketID": "1234",
                "Asset": "A12345C",
            },
        )


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_save_expect_timeout_exception(context_data, config, simple_logger):
    # GIVEN a timeout is thrown on the NATS request
    context_data.__nc__.request = AsyncMock(side_effect=TimeoutError)

    # WHEN the save method is called
    with mock.patch.object(simple_logger, "error") as mock_logger_error:
        await context_data.save(
            "test",
            {
                "Time": datetime.utcnow().strftime(DATETIME_FORMAT),
                "Result": "Tested",
                "TicketID": "1234",
                "Asset": "A12345C",
            },
        )

        # THEN expect an error log is written with the given prompt
        mock_logger_error.assert_called_with("Error saving data: request timed out")


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_find_expect_ok(context_data, config):
    # GIVEN a collection name and fake response
    coll = "test_collection"
    collection_mock = MagicMock(Collection)
    collection_mock.find = MagicMock(return_value=["response1", "response2"])
    context_data.__mongo_conn__ = {f"{config.mongo_data_db_name}": {coll: collection_mock}}

    # WHEN the find method is called
    response = await context_data.find(
        coll,
        {
            "TicketID": "1234",
            "Asset": "A5678",
        },
    )

    # THEN assert the response contains the two mocked elements
    assert len(response) == 2
    # AND assert the first element is the expected one
    assert response[0] == "response1"
    # AND assert the second element is the expected one
    assert response[1] == "response2"


@pytest.mark.unittest
@pytest.mark.parametrize(
    "collection_name, query",
    [
        (
            1234,
            {
                "TicketID": "1234",
                "Asset": "A5678",
            },
        ),
        ("collection_test", "some wrong query string"),
    ],
)
@pytest.mark.asyncio
async def test_find_wrong_type_expect_exception(context_data, config, collection_name, query):
    # GIVEN a wrong collection_name or query types
    # WHEN the find method is called
    # THEN expect an Exception to be thrown
    with pytest.raises(Exception):
        await context_data.find(collection_name, query)


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_find_expect_exception(context_data, config):
    # GIVEN the find method throws an exception
    coll = "test_collection"
    collection_mock = MagicMock(Collection)
    collection_mock.find = MagicMock(side_effect=Exception)
    context_data.__mongo_conn__ = {f"{config.mongo_data_db_name}": {coll: collection_mock}}

    # WHEN the find method is called
    # THEN expect an Exception to be thrown
    with pytest.raises(Exception):
        await context_data.find(
            coll,
            {
                "TicketID": "1234",
                "Asset": "A5678",
            },
        )
