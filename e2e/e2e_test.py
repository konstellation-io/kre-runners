import logging
import pytest
import grpc

from public_input_pb2 import Request, Response
import public_input_pb2_grpc

logger = logging.getLogger(__name__)


@pytest.mark.asyncio
async def test_main() -> None:

    with grpc.insecure_channel('localhost:9000') as channel:
        stub = public_input_pb2_grpc.EntrypointStub(channel)
        logger.info("------ Sending request to the entrypoint ------")
        response = stub.Greet(Request(name="John Doe"), timeout=5000)
        logger.info(f"Response: {response}")

    logger.info(response.greeting)
    assert response.greeting is not None
