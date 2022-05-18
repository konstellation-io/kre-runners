from google.protobuf.any_pb2 import Any

from public_input_pb2 import NodeBRequest, Response


def init(ctx):
    print("[worker init]")
    ctx.set("greeting", "Hello")


async def handler(ctx, data: Any):
    ctx.logger.info("[worker handler]")

    req = NodeBRequest()
    data.Unpack(req)

    ctx.logger.info(data)

    result = f"{req.lastname}, how are you?"
    ctx.logger.info(f"result -> {result}")

    output = Response()
    output.greeting = result
    return output
