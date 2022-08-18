#!/bin/sh

protoc -I=proto \
  --go_out=kre-go \
  --go_opt=paths=source_relative \
  --python_out=kre-py/src \
  --python_out=kre-entrypoint/src \
	--python_out=e2e/test_assets \
  proto/kre_nats_msg.proto

echo "done"
