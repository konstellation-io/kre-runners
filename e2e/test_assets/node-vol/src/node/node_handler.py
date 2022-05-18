from google.protobuf.any_pb2 import Any

from public_input_pb2 import Request, Response


def init(ctx):
    print("[worker init]")
    ctx.set("greeting", "Hello")


async def handler(ctx, data: Any):
    ctx.logger.info("[worker handler]")

    req = Request()
    data.Unpack(req)

    result = f"{ctx.get('greeting')} {req.name}!"
    ctx.logger.info(f"result -> {result}")

    output = Response()
    output.greeting = result
    return output
