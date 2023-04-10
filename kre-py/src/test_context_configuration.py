from dataclasses import dataclass
from logging import Logger
from typing import Optional
from unittest.mock import AsyncMock, MagicMock, patch

import pytest
from nats.aio.client import Client as NatsClient
from nats.js.client import JetStreamContext
from nats.js.kv import KeyValue

from context_configuration import ContextConfiguration, Scope


@dataclass
class Entry:
    bucket: str
    key: str
    value: Optional[bytes]


@pytest.fixture()
def nats_mock() -> MagicMock:
    nats_mock = MagicMock(NatsClient)

    return nats_mock


@pytest.fixture()
def jetstream_mock() -> MagicMock:
    jetstream_mock = MagicMock(JetStreamContext)

    return jetstream_mock


@pytest.fixture()
def kv_map_mock() -> dict[Scope, KeyValue]:
    kv_map = dict[Scope, KeyValue](
        {
            Scope.PROJECT: MagicMock(KeyValue),
            Scope.WORKFLOW: MagicMock(KeyValue),
            Scope.NODE: MagicMock(KeyValue),
        }
    )

    return kv_map


@pytest.fixture()
def context_configuration(
    config: MagicMock, simple_logger: Logger, kv_map_mock: dict[Scope, KeyValue]
) -> ContextConfiguration:
    context_configuration = ContextConfiguration(config, simple_logger, kv_map_mock)
    return context_configuration


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_set_configuration_default_scope_expect_ok(  # type: ignore
    context_configuration: ContextConfiguration, kv_map_mock: MagicMock
):
    # GIVEN a configuration to set
    key = "test_key"
    value = "test_value"
    kv_map_mock[Scope.NODE].put = AsyncMock()

    # WHEN the set method is called with no scope
    await context_configuration.set(key, value)

    # THEN expect the configuration to be set in default scope
    kv_map_mock[Scope.NODE].put.assert_called_with(key, bytes(value, "utf-8"))


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_set_configuration_project_scope_expect_ok(  # type: ignore
    context_configuration: ContextConfiguration, kv_map_mock: MagicMock
):
    # GIVEN a configuration to set
    key = "test_key"
    value = "test_value"
    kv_map_mock[Scope.PROJECT].put = AsyncMock()

    # WHEN the set method is called with project scope
    await context_configuration.set(key, value, Scope.PROJECT)

    # THEN expect the configuration to be set in project scope
    kv_map_mock[Scope.PROJECT].put.assert_called_with(key, bytes(value, "utf-8"))


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_set_configuration_workflow_scope_expect_ok(  # type: ignore
    context_configuration: ContextConfiguration, kv_map_mock: MagicMock
):
    # GIVEN a configuration to set
    key = "test_key"
    value = "test_value"
    kv_map_mock[Scope.WORKFLOW].put = AsyncMock()

    # WHEN the set method is called with workflow scope
    await context_configuration.set(key, value, Scope.WORKFLOW)

    # THEN expect the configuration to be set in workflow scope
    kv_map_mock[Scope.WORKFLOW].put.assert_called_with(key, bytes(value, "utf-8"))


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_set_configuration_node_scope_expect_ok(  # type: ignore
    context_configuration: ContextConfiguration, kv_map_mock: MagicMock
):
    # GIVEN a configuration to set
    key = "test_key"
    value = "test_value"
    kv_map_mock[Scope.NODE].put = AsyncMock()

    # WHEN the set method is called with node scope
    await context_configuration.set(key, value, Scope.NODE)

    # THEN expect the configuration to be set in node scope
    kv_map_mock[Scope.NODE].put.assert_called_with(key, bytes(value, "utf-8"))


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_get_configuration_default_scope_expect_ok(  # type: ignore
    context_configuration: ContextConfiguration, kv_map_mock: MagicMock
):
    # GIVEN a configuration to get
    key = "test_key"
    kv_map_mock[Scope.NODE].get = AsyncMock(
        return_value=Entry(key=key, value=None, bucket="node_bucket")
    )
    kv_map_mock[Scope.WORKFLOW].get = AsyncMock(
        return_value=Entry(key=key, value=None, bucket="workflow_bucket")
    )
    kv_map_mock[Scope.PROJECT].get = AsyncMock(
        return_value=Entry(key=key, value=bytes("test_value", "utf-8"), bucket="project_bucket")
    )

    # WHEN the get method is called with no scope
    result = await context_configuration.get(key)

    # THEN expect the configuration search for value in all scopes
    # Until it finaly encounters it in the project scope
    kv_map_mock[Scope.NODE].get.assert_called_with(key)
    kv_map_mock[Scope.WORKFLOW].get.assert_called_with(key)
    kv_map_mock[Scope.PROJECT].get.assert_called_with(key)

    assert result == "test_value"


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_get_configuration_project_scope_expect_ok(  # type: ignore
    context_configuration: ContextConfiguration, kv_map_mock: MagicMock
):
    # GIVEN a configuration to get
    key = "test_key"
    kv_map_mock[Scope.PROJECT].get = AsyncMock(
        return_value=Entry(key=key, value=bytes("test_value", "utf-8"), bucket="project_bucket")
    )

    # WHEN the get method is called with project scope
    result = await context_configuration.get(key, Scope.PROJECT)

    # THEN expect the configuration to be get in project scope
    kv_map_mock[Scope.PROJECT].get.assert_called_with(key)

    assert result == "test_value"


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_get_configuration_workflow_scope_expect_ok(  # type: ignore
    context_configuration: ContextConfiguration, kv_map_mock: MagicMock
):
    # GIVEN a configuration to get
    key = "test_key"
    kv_map_mock[Scope.WORKFLOW].get = AsyncMock(
        return_value=Entry(key=key, value=bytes("test_value", "utf-8"), bucket="workflow_bucket")
    )

    # WHEN the get method is called with workflow scope
    result = await context_configuration.get(key, Scope.WORKFLOW)

    # THEN expect the configuration to be get in workflow scope
    kv_map_mock[Scope.WORKFLOW].get.assert_called_with(key)

    assert result == "test_value"


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_get_configuration_node_scope_expect_ok(  # type: ignore
    context_configuration: ContextConfiguration, kv_map_mock: MagicMock
):
    # GIVEN a configuration to get
    key = "test_key"
    kv_map_mock[Scope.NODE].get = AsyncMock(
        return_value=Entry(key=key, value=bytes("test_value", "utf-8"), bucket="node_bucket")
    )

    # WHEN the get method is called with node scope
    result = await context_configuration.get(key, Scope.NODE)

    # THEN expect the configuration to be get in node scope
    kv_map_mock[Scope.NODE].get.assert_called_with(key)

    assert result == "test_value"


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_get_configuration_default_scope_expect_not_found(  # type: ignore
    context_configuration: ContextConfiguration, kv_map_mock: MagicMock
):
    # GIVEN a configuration to get: no value added in any scope
    key = "test_key"
    kv_map_mock[Scope.NODE].get = AsyncMock(
        return_value=Entry(key=key, value=None, bucket="node_bucket")
    )
    kv_map_mock[Scope.WORKFLOW].get = AsyncMock(
        return_value=Entry(key=key, value=None, bucket="workflow_bucket")
    )
    kv_map_mock[Scope.PROJECT].get = AsyncMock(
        return_value=Entry(key=key, value=None, bucket="project_bucket")
    )

    # WHEN the get method is called with no scope
    with patch.object(context_configuration.__logger__, "error") as mock_logger_error:
        with pytest.raises(Exception) as exception:
            await context_configuration.get(key)

    # THEN expect the configuration search for value in all scopes
    kv_map_mock[Scope.NODE].get.assert_called_with(key)
    kv_map_mock[Scope.WORKFLOW].get.assert_called_with(key)
    kv_map_mock[Scope.PROJECT].get.assert_called_with(key)

    # Until it finaly doesn't encounter it and raises an exception
    assert exception.match("No value found")
    mock_logger_error.assert_called_with(
        f"Error while getting the value for key test_key in all scopes: No value found"
    )


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_get_configuration_project_scope_expect_not_found(  # type: ignore
    context_configuration: ContextConfiguration, kv_map_mock: MagicMock
):
    # GIVEN a configuration to get: a key with no value
    key = "test_key"
    kv_map_mock[Scope.PROJECT].get = AsyncMock(
        return_value=Entry(key=key, value=None, bucket="project_bucket")
    )

    # WHEN the get method is called with project scope
    with patch.object(context_configuration.__logger__, "error") as mock_logger_error:
        with pytest.raises(Exception) as exception:
            await context_configuration.get(key, Scope.PROJECT)

    # THEN expect the configuration to try and get the value in project scope
    kv_map_mock[Scope.PROJECT].get.assert_called_with(key)

    # It doesn't encounter it and raises an exception
    assert exception.match("No value found")
    mock_logger_error.assert_called_with(
        f"Error while getting the value for key test_key: No value found"
    )


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_delete_configuration_default_scope_expect_ok(  # type: ignore
    context_configuration: ContextConfiguration, kv_map_mock: MagicMock
):
    # GIVEN a configuration to delete
    key = "test_key"
    kv_map_mock[Scope.NODE].delete = AsyncMock()

    # WHEN the delete method is called with no scope
    await context_configuration.delete(key)

    # THEN expect the configuration to be deleted in default node scope
    kv_map_mock[Scope.NODE].delete.assert_called_with(key)


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_delete_configuration_project_scope_expect_ok(  # type: ignore
    context_configuration: ContextConfiguration, kv_map_mock: MagicMock
):
    # GIVEN a configuration to delete
    key = "test_key"
    kv_map_mock[Scope.PROJECT].delete = AsyncMock()

    # WHEN the delete method is called with no scope
    await context_configuration.delete(key, Scope.PROJECT)

    # THEN expect the configuration to be deleted in project scope
    kv_map_mock[Scope.PROJECT].delete.assert_called_with(key)


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_delete_configuration_workflow_scope_expect_ok(  # type: ignore
    context_configuration: ContextConfiguration, kv_map_mock: MagicMock
):
    # GIVEN a configuration to delete
    key = "test_key"
    kv_map_mock[Scope.WORKFLOW].delete = AsyncMock()

    # WHEN the delete method is called with no scope
    await context_configuration.delete(key, Scope.WORKFLOW)

    # THEN expect the configuration to be deleted in workflow scope
    kv_map_mock[Scope.WORKFLOW].delete.assert_called_with(key)


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_delete_configuration_node_scope_expect_ok(  # type: ignore
    context_configuration: ContextConfiguration, kv_map_mock: MagicMock
):
    # GIVEN a configuration to delete
    key = "test_key"
    kv_map_mock[Scope.NODE].delete = AsyncMock()

    # WHEN the delete method is called with no scope
    await context_configuration.delete(key, Scope.NODE)

    # THEN expect the configuration to be deleted in node scope
    kv_map_mock[Scope.NODE].delete.assert_called_with(key)


@pytest.mark.unittest
@pytest.mark.asyncio
async def test_delete_configuration_default_scope_expect_not_found(  # type: ignore
    context_configuration: ContextConfiguration, kv_map_mock: MagicMock
):
    # GIVEN a configuration to delete
    key = "test_key"
    kv_map_mock[Scope.NODE].delete = AsyncMock(side_effect=Exception("Not found"))

    # WHEN the delete method is called with no scope
    # THEN expect the configuration to try to delete key in node scope
    # AND an exception is raised
    with patch.object(context_configuration.__logger__, "error") as mock_logger_error:
        with pytest.raises(Exception) as exception:
            await context_configuration.delete(key)

    kv_map_mock[Scope.NODE].delete.assert_called_with(key)
    assert exception.match("Not found")
    mock_logger_error.assert_called_with(
        f"Error while deleting the value for key test_key: Not found"
    )
