#!/bin/bash

python -m ensurepip --upgrade
pip  install asyncio nats-py
python /src/create_nats_stream.py
