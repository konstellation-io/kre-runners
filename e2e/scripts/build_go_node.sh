#!/bin/bash
# must be executed from e2e folder
  set -eu

  cd ./test_assets/nodeC-vol/src/node
  echo "Building nodeC Golang binary..."
  go build -o ../../bin/nodeC .

  cd ./../../../exitpoint-vol/src/node
  echo "Building exitpoint Golang binary..."
  go build -o ../../bin/exitpoint .

  echo "Done"
