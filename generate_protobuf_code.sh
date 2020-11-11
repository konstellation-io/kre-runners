#!/bin/sh

protoc -I=proto --go_out=kre-go --python_out=kre-py/src --python_out=kre-entrypoint/src proto/kre_nats_msg.proto

echo "done"
