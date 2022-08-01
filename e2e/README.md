# END-TO-END TESTING (E2E)

This section describes how to run the end-to-end test suite, which performs different requests to a
full local environment consisting of an entrypoint, 2 python nodes and 1 golang node.

The test launches several concurrent clients, each client will make several requests to the entrypoint.
Our point is to prove that several concurrent clients calling the entrypoint will receive all expected
responses without them being disorganized or lost.

## How to use

In order to execute the E2E test, run the following commands in two different terminals:

```sh
make build
```

```sh
make test
```

## Proto files

It's also possible to generate the protobuf files from the source code with the following command:

```sh
make proto
```
