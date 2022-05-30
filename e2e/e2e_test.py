import logging
import pytest
import grpc
import time
import random

from multiprocessing import Process, Manager
from public_input_pb2 import Request, Response
import public_input_pb2_grpc

logger = logging.getLogger(__name__)

request_per_client = 100
total_clients = 2

def assert_secuential_call(fork: int, return_list: []) -> []:
    total_requests = 0

    with grpc.insecure_channel('localhost:9000') as channel:
        stub = public_input_pb2_grpc.EntrypointStub(channel)
        logger.info(f"{fork} CLIENT STARTED ------")
        
        for i in range(request_per_client):
            time.sleep(0.01)
            name = "John Doe " + str(i) #+ str(random.randint(0, 1000000))
            request = Request(name=name)
            response = stub.Greet(request, timeout=5000)

            logger.info(f"# {fork} - {i} - Response: {response.greeting} - NAME {name}")
            
            if name not in response.greeting:
                return_list.append(f"{name} not in {response.greeting}")
                logger.info(f"{fork} ERROR - {name} not in {response.greeting}")
                assert name in response.greeting

        logger.info(f"{fork} - ---- DONE SENDING {request_per_client} REQUESTS ----")
        return return_list

#@pytest.mark.asyncio
#async def test_main() -> None:
#    assert_secuential_call()

@pytest.mark.asyncio
async def test_pararell_main() -> None:
    procs = []
    return_dict = Manager().dict()

    for client in range(total_clients):
        p = Process(target=assert_secuential_call, args=(client, return_dict))
        p.start()
        procs.append(p)

    [p.join() for p in procs]

    assert len(return_dict) is 0

    logger.info(f"END PROCESS")
