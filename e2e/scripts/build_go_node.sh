#!/bin/bash
# must be executed from e2e folder
  set -eu

  cd ./test_assets/nodeC-vol/src/node
  echo "Building nodeC Golang binary..."
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../bin/nodeC .

  cd ./../../../exitpoint-vol/src/node
  echo "Building exitpoint Golang binary..."
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../bin/exitpoint .

  echo "Done"
