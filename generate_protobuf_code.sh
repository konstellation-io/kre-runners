#!/bin/sh

protoc -I=proto \
  --go_out=kre-go \
  --go_out=go-sdk/protos \
  --go_opt=paths=source_relative \
  --python_out=kre-py/src \
  --python_out=kre-entrypoint/src proto/kai_nats_msg.proto

echo "done"
