# KRE Entrypoint 

This is the entrypoint implementation for the KRE project.

## Requirements

In order to use this entrypoint, the following environment variables must be set:

| Name                   | Description                                               |
|------------------------|-----------------------------------------------------------|
| KRT_RUNTIME_ID         | ID of the current runtime                                 |
| KRT_VERSION_ID         | ID of the current version                                 |
| KRT_VERSION            | Name of the current version                               |
| KRT_NODE_NAME          | Name of the current node                                  |
| KRT_NATS_SERVER        | NATS server URL                                           |
| KRT_NATS_SUBJECTS_FILE | Path to the file that contains the first node subject and |
| KRT_INFLUX_URI         | Influx database URI                                       |

## How it works

The entrypoint is the KRE piece that is in charge of receiving messages from gRPC, then sending the requests to the first node of the requested workflow and finally responding completed requests
back to the gRPC clients.

When the entrypoint is built, a new connection to the NATS server is established through a Jetstream's NATS client. After that, the entrypoint creates one stream per workflow specified in the version being deployed with the following format:

- `<KRT_RUNTIME_ID>-<KRT_VERSION>-<WORKFLOW_SERVICE>`

Also, the subjects to where the rest of the nodes will be subscribed will follow this format:

- `<RUNTIME_ID>-<KRT_VERSION>-<KRT_WORKFLOW_NAME>-<KRT_NODE_NAME>`

Once the stream is created, the entrypoint stores it with its subjects in a dictionary, so the
entrypoint knows all the streams and subjects that are available.

After that, the entrypoint starts to listen to a gRPC server. When a new message is received, a unique ID is generated and stored alongside the requesting gRPC client in a dictionary, so it can
be used back to respond the client.
Then an ephemeral subscription linked to the created request's ID is created to a new generated entrypoint subject (we use one subscription per message) and the message is published on the next node's subject.

At this point, subscriptions will wait until a new NATS message is published on the entrypoint's subject and once it has being received, the entrypoint will check that the message's ID is the same as the one the subscription was waiting for.
If so, the request's result is sent back to the gRPC client and the subscription is closed (Also, the NATS message is acknowledged, and the dictionary entry is removed).
