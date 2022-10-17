from google.protobuf.any_pb2 import Any

from public_input_pb2 import Request, NodeBRequest


def init(ctx):
    print("[worker init]")
    ctx.set("greeting", "Hello")


async def default_handler(ctx, data: Any) -> NodeBRequest:

    """
    This is the entrypoint handler for the nodeA workflow.

    :param ctx: The context object for the nodeA workflow.
    :param data: The message received from the previous node.

    :return: The response message to be sent to the next node.
    """

    ctx.logger.info("[worker handler]")

    req = Request()
    data.Unpack(req)

    result = f"{ctx.get('greeting')} {req.name}! greetings from nodeA"
    ctx.logger.info(f"result -> {result}")

    output = NodeBRequest()
    output.greeting = result

    await ctx.send_output(output)
