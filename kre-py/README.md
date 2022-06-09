# KRE Python Runner

This is an implementation in Python for the KRE runner.

## CUDA 10.2 Support

This image is built on top of `nvidia/cuda-10.2-devel` to add GPU support on the runner.

## How it works

The python runner is one of the two types of runners that can be used in a KRE workflow and allows executing python code.

Once the python runner is deployed, it connects to NATS and subscribes permanently to an input subject. 
Each node knows to which subject it has to subscribe and also to which subject it has to send messages, 
since the K8s manager (REFERENCE TO k8S MANAGER) tells it with environment variables. 
It's important to note that the nodes use a queue subscription, 
which allows load balancing of messages when there are multiple replicas of the runner.

When a new message is published in the input subject of a node, the runner processes the message 
and passes it to the handler, along with a context object formed by variables and useful methods for processing data. 
This handler is in charge of processing the message and returning the response to the node, 
which transforms the response to a NATS format and publishes it to the subject that corresponds to (which is indicated by an environment variable).
After that, the node ACKs the message manually.

## Usage

The injected code must implement a `handler(ctx, data)` function and optionally a `init(ctx)` function.

The context object received by these functions, has the following methods:

```python
ctx.path("relative/path.xxx")

ctx.get("label")
ctx.set("label", value)

await ctx.db.save(
  coll,
  data
)

await ctx.db.find(
  coll,
  query
)

ctx.measurement.save(
  measurement: str,
  fields: dict,
  tags: dict,
  time: datetime=None,
  precision=PRECISION_NS
)

await ctx.prediction.save(
  predicted_value: str = "",
  true_value: str = "",
  extra: dict = None,
  utcdate: datetime = None,
  error: str = ""
)
```

It is necessary to set the following environment variables in order to use the runner:

| Name                  | Description                                                         | 
|-----------------------|---------------------------------------------------------------------|
| KRT_WORKFLOW_NAME     | Name of the current workflow                                        | 
| KRT_VERSION_ID        | ID of the current version                                           | 
| KRT_VERSION           | Name of the current version                                         | 
| KRT_NODE_NAME         | Name of the current node                                            |
| KRT_NATS_SERVER       | NATS server URL                                                     |
| KRT_NATS_INPUT        | Input NATS subject to which the node will be subscribed             |
| KRT_NATS_OUTPUT       | Output NATS subject to which the node will publish the next message |
| KRT_NATS_STREAM       | NATS stream name                                                    |
| KRT_NATS_MONGO_WRITER | Mongo writer name                                                   |
| KRT_BASE_PATH         | Base path where the src folder is located                           |
| KRT_HANDLER_PATH      | Path to the handler file                                            |
| KRT_MONGO_URI         | Mongo database URI                                                  |
| KRT_INFLUX_URI        | Influx database URI                                                 |

This is an example of the code that will be run by the docker py3 runner:

```python
import pickle
import numpy as np
import pandas as pd

# Import the proto message types
from public_input_pb2 import Request, Response

# this function will be executed once when the runner is starting
def init(ctx):
  # load file and save in memory to be used within the handler
  ctx.set("categories", pickle.load(ctx.path("data/categories.pkl")))

# this function will be executed when a message is received
async def handler(ctx, data):
  categories = ctx.get("categories")

  # data is the received message from the queue
  req = Request()
  data.Unpack(req)

  normalized_data = np.xxx(categories)
  normalized_data = pd.xxx(normalized_data)

  res = Response()
  res.normalized_data = normalized_data

  return res # return a protobuf for the next node
```

## Development

First of all, install the dependencies using pipenv:

```bash
pipenv install --dev
```

If you don't have pipenv installed (you must have python 3.7 installed in your system):

```shell script
pip3 install --user pipenv
```

### Integration tests

You can run the integration tests following these steps:

1. Start the NATS, MongoDB and InfluxDB services:

```shell script
cd test
docker-compose up
```

2. When the services are ready then run the integration tests:

```shell script
pipenv shell
PYTHONPATH=src pytest -vv -m integration
```
