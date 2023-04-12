from google.protobuf.any_pb2 import Any

from public_input_pb2 import NodeBRequest, NodeCRequest


def init(ctx) -> None:
    print("[worker init]")


async def default_handler(ctx, data: Any) -> None:

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
    output.testing.CopyFrom(req.testing)

    await ctx.send_output(output)
