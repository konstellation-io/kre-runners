#!/bin/sh
# must be executed from e2e folder

# For Go testing utility
protoc -I=./test_assets/entrypoint-vol/krt-files \
  --go_out=test_go \
  --go-grpc_out=test_go \
  ./test_assets/entrypoint-vol/krt-files/public_input.proto

# For Go node
protoc -I=test_assets/entrypoint-vol/krt-files/ \
  --go_out=test_assets/nodeC-vol/src/node \
  --go_out=test_assets/exitpoint-vol/src/node \
  test_assets/entrypoint-vol/krt-files/public_input.proto

# For Python nodes and entrypoint
protoc -I=test_assets/entrypoint-vol/krt-files/ \
  --python_out=test_assets/nodeA-vol/src/node \
  --python_out=test_assets/nodeB-vol/src/node \
  test_assets/entrypoint-vol/krt-files/public_input.proto

  echo "Done"
