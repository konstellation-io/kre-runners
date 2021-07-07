#!/bin/bash

echo "Generating public_input.proto code..."

protoc -I=test_assets \
       --python_out=src/test_utils \
       --python_out=test_assets/myvol/src/node \
       test_assets/public_input_for_testing.proto

echo "Done"
