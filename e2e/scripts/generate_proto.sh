#!/bin/sh

# For Go testing utility
protoc -I=./test_assets/entrypoint-vol/krt-files \
  --go_out=test_go \
  --go-grpc_out=test_go \
  ./test_assets/entrypoint-vol/krt-files/public_input.proto

# For Go node
protoc -I=test_assets/entrypoint-vol/krt-files/ \
  --go_out=test_assets/nodeC-vol/src/node \
  test_assets/entrypoint-vol/krt-files/public_input.proto

python3 -m grpc_tools.protoc \
  -I=test_assets/entrypoint-vol/krt-files/ \
  --python_out=. \
  --grpc_python_out=. \
  test_assets/entrypoint-vol/krt-files/public_input.proto

  echo "Done"
