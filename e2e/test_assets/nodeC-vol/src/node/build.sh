#!/bin/bash

# shellcheck disable=SC2086

## USAGE:
#    ./build.sh 

set -eu

echo "Building nodeC Golang binary..."
go build -o ../../bin/nodeC .
