import time
from google.protobuf.any_pb2 import Any

from public_input_pb2 import Request, NodeBRequest, Response


async def init(ctx) -> None:
    print("[worker init]")
    await ctx.configuration.set("greeting", "Hello")


async def default_handler(ctx, data: Any) -> None:
    """
    This is the entrypoint handler for the nodeA workflow.

    :param ctx: The context object for the nodeA workflow.
    :param data: The message received from the previous node.

    :return: The response message to be sent to the next node.
    """

    ctx.logger.info("[worker handler]")

    req = Request()
    data.Unpack(req)

    if req.testing.is_early_reply:
        ctx.logger.info(f"early reply recieved")
        output_er = Response()
        output_er.greeting = req.name
        await ctx.send_early_exit(output_er)
        return

    if req.testing.is_early_exit:
        ctx.logger.info(f"early exit recieved")
        output_ee = Response()
        output_ee.greeting = req.name
        await ctx.send_early_reply(output_ee)
        time.sleep(0.3)  # let exitpoint finish its requests

    greeting = await ctx.configuration.get("greeting")
    ctx.logger.info(f"greeting -> {greeting}")
    result = f"{greeting} {req.name}! greetings from nodeA"
    ctx.logger.info(f"result -> {result}")

    output = NodeBRequest()
    output.greeting = result
    output.testing.CopyFrom(req.testing)

    await ctx.send_output(output)
