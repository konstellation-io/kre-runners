import asyncio
import os
from datetime import datetime
import logging
import time
from multiprocessing import Process
from unittest import mock

import pytest
from nats.aio.client import Client as NATS
from nats.js import JetStreamContext
from nats.js.api import DeliverPolicy, ConsumerConfig

from config import Config
from kre_nats_msg_pb2 import KreNatsMessage
from main import NodeRunner
from test_utils.public_input_for_testing_pb2 import Request, Response

TEST_ENV_VARS = {
    "KRT_WORKFLOW_NAME": "workflow1",
    "KRT_VERSION_ID": "version.1234",
    "KRT_VERSION": "testVersion1",
    "KRT_NODE_NAME": "nodeA",
    "KRT_NATS_SERVER": "localhost:4222",
    "KRT_NATS_INPUT": "test-subject-input",
    "KRT_NATS_OUTPUT": "",
    "KRT_NATS_ENTRYPOINT_SUBJECT": "entrypoint_subject",
    "KRT_NATS_MONGO_WRITER": "mongo_writer",
    "KRT_BASE_PATH": f"{os.getcwd()}/test_assets/myvol",
    "KRT_HANDLER_PATH": "src/node/node_handler.py",
    "KRT_MONGO_URI": "mongodb://admin:123456@localhost:27017/admin",
    "KRT_INFLUX_URI": "http://localhost:8088",
}

logger = logging.getLogger(__name__)


@pytest.fixture(autouse=True)
def runner():
    print("Starting the node runner...")

    #with mock.patch.dict(os.environ, TEST_ENV_VARS):
        #runner = NodeRunner()

    #p = Process(target=runner.start)
    #p.start()
    # yield runner
    #p.terminate()


def prepare_nats_message():
    # Prepare request message
    req_nats_msg = KreNatsMessage()

    # Add tracking info from the previous node
    tracking = req_nats_msg.tracking.add()
    tracking.node_name = "previous_node"
    tracking.start = datetime.utcnow().isoformat()
    tracking.end = datetime.utcnow().isoformat()

    # Put the request message inside a KreNatsMessage
    req = Request()
    req.name = "John Doe"
    req_nats_msg.payload.Pack(req)

    return req_nats_msg.SerializeToString()


@pytest.mark.asyncio
async def test_main() -> None:
    input_subject = TEST_ENV_VARS["KRT_NATS_INPUT"]
    stream = TEST_ENV_VARS["KRT_VERSION_ID"].replace('.', '-') + "-" + TEST_ENV_VARS["KRT_WORKFLOW_NAME"]
    output_subject = "entrypoint_subject"

    logger.info("Connecting to NATS...")
    nc = NATS()
    await nc.connect(TEST_ENV_VARS["KRT_NATS_SERVER"])
    # Create JetStream context
    js = nc.jetstream()

    logger.info(f"Adding stream {stream}...")
    await js.add_stream(name=stream, subjects=[input_subject, output_subject])

    nats_message = prepare_nats_message()

    sub = await js.subscribe(
        stream=stream,
        subject=output_subject,
        durable=stream,
        config=ConsumerConfig(
            deliver_policy=DeliverPolicy.ALL,
        ),
    )

    time.sleep(5)
    logger.info(f"Sending a test message to {input_subject}...")
    ack = await js.publish(stream=stream, subject=input_subject, payload=nats_message)
    logger.info(ack)

    msg = await sub.next_msg(timeout=1000)
    await msg.ack()

    # Get the response
    res_nats_msg = KreNatsMessage()
    res_nats_msg.ParseFromString(msg.data)

    res = Response()
    res_nats_msg.payload.Unpack(res)

    assert res_nats_msg.error == ""
    assert res.greeting == "Hello John Doe!"

    nc.close()


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
@mock.patch(os.environ, environment_variables)
def node_runner(mocked_nats: mock.Mock, ) -> NodeRunner:
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
        "KRT_NATS_INPUT": "runtime-1-version-1-workflow-a.node-a",
        "KRT_NATS_OUTPUT": "runtime-1-version-1-workflow-a.node-b",
        "KRT_WORKFLOW_NAME": "workflow-a",
        "KRT_NATS_MONGO_WRITER": "mongo_writer",
        "KRT_BASE_PATH": "/tmp",
        "KRT_HANDLER_PATH": "src/node/node_handler.py",
        "KRT_MONGO_URI": "mongodb://mock:mock@mock:27017/admin",
    }

    with mock.patch.dict(os.environ, environment_variables):
        with mock.patch.object(NodeRunner, "load_handler", return_value=None):
            runner = NodeRunner()
            runner.nc = mocked_nats
            return runner


@pytest.mark.asyncio
async def test_create_message_cb(node_runner: NodeRunner) -> None:
    msg = "Hi! "
    response = node_runner.create_message_cb()

    node_runner.handler_fn.return_value = "Hello!"

    await response.__call__(msg)



@pytest.mark.asyncio
async def test_process_messages(node_runner: NodeRunner) -> None:
    node_runner.js = node_runner.nc.jetstream()

    await node_runner.process_messages()

    assert node_runner.js.subscribe.call_count == 1
