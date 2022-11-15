from google.protobuf.any_pb2 import Any

from public_input_pb2 import NodeBRequest, NodeCRequest


def init(ctx):
    print("[worker init]")


async def default_handler(ctx, data: Any) -> NodeCRequest:

    """
    This is the entrypoint handler for the nodeB workflow.

    :param ctx: The context object for the nodeB workflow.
    :param data: The message received from the previous node.

    :return: The response message to be sent to the next node.
    """
    ctx.logger.info("[worker handler]")

    if ctx.is_message_early_reply() or ctx.is_message_early_exit():
        ctx.logger.info("ignoring request")
        return

    req = NodeBRequest()
    data.Unpack(req)

    result = f"{req.greeting}, nodeB"
    ctx.logger.info(f"result -> {result}")

    output = NodeCRequest()
    output.greeting = result

    await ctx.send_output(output)
