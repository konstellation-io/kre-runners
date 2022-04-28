#!/bin/sh

protoc -I=proto --go_out=kre-go --go_opt=paths=source_relative --python_out=kre-py/src --python_out=kre-entrypoint/src proto/kre_nats_msg.proto

echo "done"
