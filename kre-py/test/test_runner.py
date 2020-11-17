import asyncio
import json

from nats.aio.client import Client as NATS
from nats.aio.errors import ErrTimeout


from kre_nats_msg_pb2 import KreNatsMessage
from test_msg_pb2 import TestInput, TestOutput


async def run(loop):
    print("Connecting to NATS...")
    nc = NATS()
    await nc.connect("localhost:4222", loop=loop)

    input_subject = "test-subject-input"
    try:
        req = TestInput()
        req.name = "John Doe"

        natsMsg = KreNatsMessage()
        natsMsg.payload.Pack(req)

        async def message_handler(msg):
            subject = msg.subject
            reply = msg.reply
            data = msg.data.decode()
            print(
                f"[WRITER] Received a message on '{subject} {reply}': {data}")
            try:
                await nc.publish(reply, bytes("{\"success\": true }", encoding='utf-8'))
                print(f"[WRITER] reply ok")
            except Exception as err:
                print(f"[WRITER] exception: {err}")

        sid = await nc.subscribe("mongo_writer", cb=message_handler)
        print("Waiting for a metrics message...")

        # Stop receiving after 3 calls to prediction.save() and 1 call to db.save().
        await nc.auto_unsubscribe(sid, 4)

        print(f"Sending a test message to {input_subject}...")
        msg = await nc.request(input_subject, natsMsg.SerializeToString(), timeout=10)

        res = KreNatsMessage()
        res.ParseFromString(msg.data)

        out = TestOutput()
        res.payload.Unpack(out)

        print("error -> ", res.error)
        print("greeting -> ", out.greeting)
    except ErrTimeout:
        print("Request timed out")

    await nc.close()


if __name__ == '__main__':
    loop = asyncio.get_event_loop()
    loop.run_until_complete(run(loop))
    loop.close()
