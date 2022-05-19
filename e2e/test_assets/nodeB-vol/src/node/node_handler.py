from google.protobuf.any_pb2 import Any

from public_input_pb2 import NodeBRequest, NodeCRequest


def init(ctx):
    print("[worker init]")
    ctx.set("greeting", "Hello")


async def handler(ctx, data: Any) -> NodeCRequest:

    """
    This is the entrypoint handler for the nodeB workflow.

    :param ctx: The context object for the nodeB workflow.
    :param data: The message received from the previous node.

    :return: The response message to be sent to the next node.
    """

    ctx.logger.info("[worker handler]")

    req = NodeBRequest()
    data.Unpack(req)

    result = f"{req.lastname}, how are you?"
    ctx.logger.info(f"result -> {result}")

    output = NodeCRequest()
    output.greeting = result
    return output
