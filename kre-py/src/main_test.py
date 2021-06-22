import os
from datetime import datetime
from multiprocessing import Process
from unittest import mock

import pytest
from nats.aio.client import Client as NATS

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
    "KRT_NATS_MONGO_WRITER": "mongo_writer",
    "KRT_BASE_PATH": f"{os.getcwd()}/test_assets/myvol",
    "KRT_HANDLER_PATH": "src/node/node_handler.py",
    "KRT_MONGO_URI": "mongodb://admin:123456@localhost:27017/admin",
    "KRT_INFLUX_URI": "http://localhost:8088",
}


@pytest.fixture(autouse=True)
def runner():
    print("Starting the node runner...")

    with mock.patch.dict(os.environ, TEST_ENV_VARS):
        runner = NodeRunner()

    p = Process(target=runner.start)
    p.start()
    yield runner
    p.terminate()


@pytest.mark.integration
@pytest.mark.asyncio
async def test_main() -> None:
    print("Connecting to NATS...")
    nc = NATS()
    await nc.connect(TEST_ENV_VARS["KRT_NATS_SERVER"])

    input_subject = TEST_ENV_VARS["KRT_NATS_INPUT"]

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

    print(f"Sending a test message to {input_subject}...")
    msg = await nc.request(input_subject, req_nats_msg.SerializeToString(), timeout=5)

    # Get the response
    res_nats_msg = KreNatsMessage()
    res_nats_msg.ParseFromString(msg.data)

    res = Response()
    res_nats_msg.payload.Unpack(res)

    assert res_nats_msg.error == ""
    assert res.greeting == "Hello John Doe!"

    nc.close()
