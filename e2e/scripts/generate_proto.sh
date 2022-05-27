#!/bin/sh

protoc -I=./test_assets/entrypoint-vol/krt-files \
  --go_out=test_go --go_opt=paths=source_relative\
  --go-grpc_out=test_go --go-grpc_opt=paths=source_relative \
  ./test_assets/entrypoint-vol/krt-files/*.proto

  echo "Done"