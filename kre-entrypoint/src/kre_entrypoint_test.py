import asyncio
import pytest
import mock
import logging
import os

from nats.js import JetStreamContext
from nats.aio.subscription import Subscription
from nats.aio.msg import Msg
from grpclib.server import Stream
from kre_entrypoint import EntrypointKRE
from config import Config
import kre_nats_msg_pb2

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
    environment_variables = {
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

    with mock.patch.dict(os.environ, environment_variables, clear=True):
        config = Config()
        entrypoint = EntrypointKRE(logger, {}, config)

        entrypoint.nc = mocked_nats

        return entrypoint


@pytest.mark.asyncio
async def test_streams_should_be_created_with_proper_names(entrypoint: EntrypointKRE) -> None:

    # GIVEN an entrypoint subject with name "workflow-a"
    # AND a runtime id "runtime-1"
    # AND a version id "version-1"
    entrypoint.output_subjects = {
        "workflow-a": "runtime-1-version-1-workflow-a.node-a",
    }

    entrypoint.config.krt_runtime_id = "runtime1"
    entrypoint.config.krt_version_id = "version1"

    expected_stream_name = "runtime1-version1-workflow-a"
    expected_subjects = [expected_stream_name + ".*"]

    # WHEN the entrypoint is initialized
    await entrypoint.start()

    # THEN the streams should be created with proper names
    # AND only one stream should be created for the subject "workflow-a"
    entrypoint.js.add_stream.assert_called_with(name=expected_stream_name, subjects=expected_subjects)
    assert entrypoint.js.add_stream.call_count == 1


@pytest.mark.asyncio
async def test_start_should_create_all_workflow_streams(entrypoint: EntrypointKRE) -> None:

    # GIVEN three entrypoint subjects with names "workflow-a", "workflow-b", "workflow-c"
    entrypoint.output_subjects = {
        "workflow-a": "runtime-1-version-1-workflow-a.node-a",
        "workflow-b": "runtime-1-version-1-workflow-b.node-a",
        "workflow-c": "runtime-1-version-1-workflow-b.node-a",
    }

    # WHEN the entrypoint is initialized
    await entrypoint.start()

    # THEN three streams should be created for the subjects "workflow-a", "workflow-b", "workflow-c"
    assert entrypoint.js.add_stream.call_count == 3


@pytest.mark.asyncio
async def test_process_grpc_message_should_a_configurable_timeout_upon_nats_subscriptions(
    entrypoint: EntrypointKRE
) -> None:

    # GIVEN an entrypoint subject with name "workflow-a"
    # AND a subscription
    # AND a stream
    # AND a configured nats_subscriptions_timeout of 927 seconds
    # AND a gRPC stream
    workflow_name = "workflow-a"

    message_mock = mock.Mock(spec=Msg)
    subscription_mock = mock.Mock(spec=Subscription)
    stream_mock = mock.Mock()
    grpc_stream_mock = mock.Mock(spec=Stream)

    message_mock.data = b"Hello World"
    subscription_mock.next_msg.return_value = message_mock
    grpc_stream_mock.recv_message.return_value = kre_nats_msg_pb2.KreNatsMessage()
    grpc_stream_mock.send_message.return_value = kre_nats_msg_pb2.KreNatsMessage()

    entrypoint.output_subjects = {
        workflow_name: "runtime-1-version-1-workflow-a.node-a",
    }

    entrypoint.subscriptions = {
        workflow_name: subscription_mock
    }

    entrypoint.streams = {
        workflow_name: stream_mock
    }

    entrypoint.config.request_timeout = 927
    entrypoint.js = entrypoint.nc.jetstream()

    # WHEN the message is received
    await entrypoint.process_grpc_message(grpc_stream_mock, workflow_name)

    # THEN the timeout should be configurable
    subscription_mock.next_msg.assert_called_with(timeout=entrypoint.config.request_timeout)
