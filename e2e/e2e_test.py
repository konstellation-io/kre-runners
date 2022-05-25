import logging
import pytest
import grpc
import random

from public_input_pb2 import Request, Response
import public_input_pb2_grpc

logger = logging.getLogger(__name__)


@pytest.mark.asyncio
async def test_main() -> None:

    with grpc.insecure_channel('localhost:9000') as channel:
        stub = public_input_pb2_grpc.EntrypointStub(channel)
        logger.info("------ Sending request to the entrypoint ------")

        for i in range(20):
            name = "John Doe " + str(random.randint(0, 1000000))
            request = Request(name=name)
            response = stub.Greet(request, timeout=5000)
            logger.info(f"Response: {response}")
            
            assert name in response.greeting

