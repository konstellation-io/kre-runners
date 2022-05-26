import logging
import pytest
import grpc
import time
import random

from multiprocessing import Process
from public_input_pb2 import Request, Response
import public_input_pb2_grpc

logger = logging.getLogger(__name__)

def assert_secuential_call():
    with grpc.insecure_channel('localhost:9000') as channel:
        stub = public_input_pb2_grpc.EntrypointStub(channel)
        logger.info("------ Sending request to the entrypoint ------")

        for i in range(100):
            name = "John Doe " + str(i) #+ str(random.randint(0, 1000000))
            request = Request(name=name)
            response = stub.Greet(request, timeout=5000)
            logger.info(f"Response: {response}")
            
            assert name in response.greeting

        logger.info("---- DONE SENDING 100 REQUESTS ----")

#@pytest.mark.asyncio
#async def test_main() -> None:
#    assert_secuential_call()

@pytest.mark.asyncio
async def test_pararell_main() -> None:
    run_cpu_tasks_in_parallel([assert_secuential_call, assert_secuential_call])

    
def run_cpu_tasks_in_parallel(tasks):
    running_tasks = [Process(target=task) for task in tasks]
    for running_task in running_tasks:
        running_task.start()

    time.sleep(1)

    for running_task in running_tasks:
        running_task.join()

