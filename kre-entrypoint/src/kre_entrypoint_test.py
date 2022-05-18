import asyncio
import pytest
import mock
import logging
import os

from nats.js import JetStreamContext
from kre_entrypoint import EntrypointKRE
from config import Config

logger = logging.getLogger(__name__)

@pytest.fixture
def mocked_nats() -> mock.Mock:
    mocked_nats = mock.Mock()

    future = asyncio.Future()
    future.set_result(123)
    mock_js = mock.Mock(spec=JetStreamContext)

    mocked_nats.connect.return_value = future
    mocked_nats.jetstream.return_value = mock_js
    return mocked_nats

@pytest.fixture
def entrypoint(mocked_nats: mock.Mock, ) -> EntrypointKRE:
    ENV = {
        "KRT_VERSION_ID": "v1",
        "KRT_VERSION": "version1",
        "KRT_RUNTIME_ID": "runtime1",
        "KRT_NODE_NAME": "entrypoint",
        "KRT_NODE_ID": "entrypoint-id",
        "KRT_NATS_SERVER": "mock:4222",
        "KRT_NATS_SUBJECTS_FILE": "/src/conf/subjects.json",
        "KRT_INFLUX_URI": "http://mock:8088",
        "KRT_NATS_STREAM": "runtime-1-version-1-workflow-a",
        "KRT_NATS_INPUT": "runtime-1-version-1-workflow-a.entrypoint",
        "KRT_NATS_OUTPUT": "runtime-1-version-1-workflow-a.node-a",
    }

    with mock.patch.dict(os.environ, ENV, clear=True):
        config = Config()
        entrypoint = EntrypointKRE(logger, {}, config)

        entrypoint.nc = mocked_nats

        return entrypoint


@pytest.mark.asyncio
async def test_streams_should_be_created_with_proper_names(entrypoint: EntrypointKRE) -> None:

    entrypoint.subjects = {
        "workflow-a": "runtime-1-version-1-workflow-a.node-a",
    }

    entrypoint.config.krt_runtime_id = "runtime1"
    entrypoint.config.krt_version_id = "version1"

    expected_stream_name = "runtime1-version1-workflow-a"
    expected_subjects = [expected_stream_name + ".*"]

    await entrypoint.start()

    entrypoint.js.add_stream.assert_called_with(name=expected_stream_name, subjects=expected_subjects)


@pytest.mark.asyncio
async def test_start_should_create_one_workflow_streams(entrypoint: EntrypointKRE) -> None:

    entrypoint.subjects = {
        "workflow-a": "runtime-1-version-1-workflow-a.node-a",
    }

    await entrypoint.start()

    assert entrypoint.js.add_stream.call_count == 1


@pytest.mark.asyncio
async def test_start_should_create_all_workflow_streams(entrypoint: EntrypointKRE) -> None:

    entrypoint.subjects = {
        "workflow-a": "runtime-1-version-1-workflow-a.node-a",
        "workflow-b": "runtime-1-version-1-workflow-b.node-a",
        "workflow-c": "runtime-1-version-1-workflow-b.node-a",
    }

    await entrypoint.start()

    assert entrypoint.js.add_stream.call_count == entrypoint.subjects.__len__()
