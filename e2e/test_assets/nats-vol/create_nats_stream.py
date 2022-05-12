import asyncio
import json
import os

import nats

NATS_STREAMS_CONFIG = os.environ.get("NATS_STREAMS_CONFIG")


async def main():
    # Open json containing nats streaming and subjects
    with open(NATS_STREAMS_CONFIG) as f:
        config = json.load(f)

    try:
        nc = await nats.connect("nats:4222")

        # Create JetStream context.
        js = nc.jetstream()

        print(config)
        print(config["stream"])
        print(config["subjects"])

        # Create NATS JS stream.
        await js.add_stream(name=config["stream"], subjects=config["subjects"])
    except Exception as e:
        print(f"Error while creating NATS JS stream: {e}")
        raise e


if __name__ == "__main__":
    asyncio.run(main())