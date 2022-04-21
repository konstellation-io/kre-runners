import os
from datetime import datetime

import pytest
from nats.aio.client import Client as NATS

from kre_nats_msg_pb2 import KreNatsMessage
from test_utils.public_input_for_testing_pb2 import Request, Response


@pytest.mark.integration
@pytest.mark.asyncio
async def test_main() -> None:
    print("Connecting to NATS...")
    nc = NATS()
    await nc.connect(os.getenv("KRT_NATS_SERVER", "localhost:4222"))

    input_subject = os.getenv("KRT_NATS_INPUT", "test-subject-input")

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

    msg = await nc.request(input_subject, req_nats_msg.SerializeToString(), timeout=20)

    # Get the response
    res_nats_msg = KreNatsMessage()
    res_nats_msg.ParseFromString(msg.data)

    res = Response()
    res_nats_msg.payload.Unpack(res)

    assert res_nats_msg.error == ""
    assert res.greeting == "Hello John Doe!"

    await nc.close()
