import logging
import multiprocessing

import pytest
import grpc
import time
import random

from multiprocessing import Process
from public_input_pb2 import Request
import public_input_pb2_grpc

logger = logging.getLogger(__name__)


def assert_sequential_call():
    with grpc.insecure_channel('localhost:9000') as channel:
        stub = public_input_pb2_grpc.EntrypointStub(channel)
        process = multiprocessing.current_process()
        logger.info(f"------ Sending request to the entrypoint from process {process.name} ------")

        for i in range(100):
            time.sleep(random.randint(1, 10))
            name = "John Doe " + str(i)  # + str(random.randint(0, 1000000))
            request = Request(name=name)
            response = stub.Greet(request, timeout=5000)
            # logger.info(f"Response: {response}")

            assert name in response.greeting

        logger.info(f"---- DONE SENDING 100 REQUESTS IN PROCESS {process.name}----")


@pytest.mark.asyncio
async def test_parallel_main() -> None:
    run_cpu_tasks_in_parallel([
        assert_sequential_call,
        assert_sequential_call,
        assert_sequential_call,
        assert_sequential_call,
    ])


def run_cpu_tasks_in_parallel(tasks):
    running_tasks = [Process(target=task) for task in tasks]
    for running_task in running_tasks:
        running_task.start()

    for running_task in running_tasks:
        running_task.join()
