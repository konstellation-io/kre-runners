#!/bin/bash
# must be executed from e2e folder
  set -eu

  cd ./test_go
  echo "Running e2e Go test..."
  go run main.go

  echo "Done"
