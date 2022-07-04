# KRE Golang Runner 

This is an implementation in Go for the KRE runner.

## How it works

The Go runner is one of the two types of runners that can be used in a KRE workflow and allows
executing go code.

Once the go runner is deployed, it connects to NATS and subscribes permanently to an input
subject.
Each node knows to which subject it has to subscribe and also to which subject it has to send messages,
since the [K8s manager](https://github.com/konstellation-io/kre/tree/main/engine/k8s-manager) tells it with environment variables.
It's important to note that the nodes use a queue subscription,
which allows load balancing of messages when there are multiple replicas of the runner.

When a new message is published in the input subject of a node, the runner passes it down to a
handler function, along with a context object formed by variables and useful methods for processing data.
This handler is the solution implemented by the client and given in the krt file generated.
Once executed, the result will be taken by the runner and transformed into a NATS message that
will then be published to the next node's subject (indicated by an environment variable).
After that, the node ACKs the message manually.

## Requirements

It is necessary to set the following environment variables in order to use the runner:

| Name                  | Description                                                         |
|-----------------------|---------------------------------------------------------------------|
| KRT_WORKFLOW_NAME     | Name of the current workflow                                        |
| KRT_VERSION_ID        | ID of the current version                                           |
| KRT_VERSION           | Name of the current version                                         |
| KRT_NODE_NAME         | Name of the current node                                            |
| KRT_NATS_SERVER       | NATS server URL                                                     |
| KRT_NATS_INPUT        | Input NATS subject to which the node will be subscribed             |
| KRT_NATS_OUTPUT       | Output NATS subject to which the node will publish the next message |
| KRT_NATS_STREAM       | NATS stream name                                                    |
| KRT_NATS_MONGO_WRITER | Mongo writer name                                                   |
| KRT_BASE_PATH         | Base path where the src folder is located                           |
| KRT_HANDLER_PATH      | Path to the handler file                                            |
| KRT_MONGO_URI         | Mongo database URI                                                  |
| KRT_INFLUX_URI        | Influx database URI                                                 |

## Run Tests

Execute the test running:

``` sh
go test
```

Integration tests have been migrated to the e2e folder.
