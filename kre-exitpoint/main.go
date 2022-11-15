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

	if ctx.IsMessageEarlyReply() || ctx.IsMessageEarlyExit() {
		ctx.SendOutput(data)
	}

	return nil
}

// lastNodeHandler will redirect all incoming messages from the last node to the entrypoint
func lastNodeHandler(ctx *kre.HandlerContext, data *anypb.Any) error {
	ctx.Logger.Info("[exitpoint handler invoked]")

	ctx.SendAny(data)

	return nil
}

func main() {
	logger := simplelogger.New(simplelogger.LevelInfo)
	config := config.NewConfig(logger)

	// kre will assign the default exitpoint to subscribe to all workflow nodes
	handlers := map[string]kre.Handler{
		config.LastNodeName: lastNodeHandler,
	}

	kre.Start(handlerInit, defaultHandler, handlers)
}
