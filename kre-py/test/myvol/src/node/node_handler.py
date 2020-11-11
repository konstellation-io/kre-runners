from test_msg_pb2 import TestInput, TestOutput


def init(ctx):
    print("[worker init]")
    ctx.set("greeting", "Hello")


async def handler(ctx, data):
    ctx.logger.info("[worker handler]")

    req = TestInput()
    data.Unpack(req)

    result = f"{ctx.get('greeting')} {req.name}!"
    ctx.logger.info(f"result -> {result}")

    ctx.logger.info("Saving some metrics...")
    # Saves prediction metrics in MongoDB DB sending a message to the MongoWriter queue
    await ctx.prediction.save(predicted_value="class_x", true_value="class_y")
    await ctx.prediction.save(error=ctx.prediction.ERR_MISSING_VALUES)
    await ctx.prediction.save(
        error=ctx.prediction.ERR_NEW_LABELS)  # If the date is not set, the 'date' field value will be now

    ctx.logger.info("Saving some data...")
    # Saves data in MongoDB DB sending a message to the MongoWriter queue
    await ctx.db.save("test_data", {"random": 123, "dict": {"some": True, "value": 0}})

    output = TestOutput()
    output.greeting = result
    return output
