#!/bin/sh

protoc -I=./test_assets/entrypoint-vol/krt-files \
  --go_out=test_go \
  --go-grpc_out=test_go \
  ./test_assets/entrypoint-vol/krt-files/*.proto

  echo "Done"
