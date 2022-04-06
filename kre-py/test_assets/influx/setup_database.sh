#!/bin/bash

sleep 5

influx -host influx -port 8086 -execute "CREATE DATABASE ${INFLUX_DATABASE}" && echo "DONE"