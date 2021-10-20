# KRE Python Runner

This is an implementation in Python for the KRE runner.

## CUDA 10.2 Support

This image is built on top of `nvidia/cuda-10.2-devel` to add GPU support on the runner.

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

The runner will have the following environment variables:

```bash
KRT_WORKFLOW_NAME
KRT_VERSION_ID
KRT_VERSION
KRT_NODE_NAME
KRT_NATS_SERVER
KRT_NATS_INPUT
KRT_NATS_OUTPUT
KRT_NATS_MONGO_WRITER
KRT_BASE_PATH
KRT_HANDLER_PATH
KRT_MONGO_URI
KRT_INFLUX_URI
```

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
