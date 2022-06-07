# END-TO-END TESTING (E2E)

This section describes how to run the end-to-end test which performs different request to a full local environment with an entrypoint, 2 python nodes and 1 golang node.

The test checks if the requests are properly responded and also responded in the same order as the requests were sent.

## How to use

In order to execute the E2E test, run the following commands in two different terminals:

```sh
$ make build
```

```sh
$ make test
```

## Proto files

It's also possible to generate the protobuf files from the source code with the following command:

```sh
$ make proto
```
