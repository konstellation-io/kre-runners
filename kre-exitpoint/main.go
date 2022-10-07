package main

import (
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/konstellation-io/kre-runners/exitpoint/config"
	"github.com/konstellation-io/kre-runners/kre-go"
	"github.com/konstellation-io/kre/libs/simplelogger"
)

func handlerInit(ctx *kre.HandlerContext) {
	ctx.Logger.Info("[exitpoint init]")
}

// defaultHandler will redirect only early reply and early exit messages to the entrypoint
func defaultHandler(ctx *kre.HandlerContext, data *anypb.Any) error {
	ctx.Logger.Info("[exitpoint default handler invoked]")

	msgType := ctx.GetRequestMessageType()
	if msgType == kre.MessageType_EARLY_REPLY || msgType == kre.MessageType_EARLY_EXIT {
		ctx.SendOutput(data)
	}

	saveExitpointMetrics(ctx)

	return nil
}

// lastNodeHandler will redirect all incoming messages from the last node to the entrypoint
func lastNodeHandler(ctx *kre.HandlerContext, data *anypb.Any) error {
	ctx.Logger.Info("[exitpoint handler invoked]")

	ctx.SendAny(data)

	saveExitpointMetrics(ctx)

	return nil
}

// saveExitpointMetrics is a helper function used to save influxdb metrics.
// We will save one metric: One to count the number of times this node has been called.
func saveExitpointMetrics(ctx *kre.HandlerContext) {
	tags := map[string]string{}

	fields := map[string]interface{}{
		"called_node": "exitpoint",
	}

	ctx.Measurement.Save("number_of_calls", fields, tags)

}

func main() {
	logger := simplelogger.New(simplelogger.LevelInfo)
	config := config.NewConfig(logger)

	handlers := map[string]kre.Handler{
		config.LastNodeName: lastNodeHandler,
	}

	kre.Start(handlerInit, defaultHandler, handlers)
}
