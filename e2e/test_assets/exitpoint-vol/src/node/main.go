package main

import (
	"exitpoint/proto"

	"google.golang.org/protobuf/types/known/anypb"

	"github.com/konstellation-io/kre-runners/kre-go/v4"
)

var attendedRequests map[string]bool

func handlerInit(ctx *kre.HandlerContext) {
	ctx.Logger.Info("[worker init]")
	attendedRequests = make(map[string]bool)
}

func defaultHandler(ctx *kre.HandlerContext, data *anypb.Any) error {
	ctx.Logger.Info("[worker handler default]")

	return nil
}

func nodeA_handler(ctx *kre.HandlerContext, data *anypb.Any) error {
	ctx.Logger.Info("[worker handler for nodeA]")

	if ctx.IsMessageEarlyReply() || ctx.IsMessageEarlyExit() {
		ctx.SendAny(data)
	}
	if ctx.IsMessageEarlyReply() {
		attendedRequests[ctx.GetRequestID()] = true
	}

	return nil
}

func nodeC_handler(ctx *kre.HandlerContext, data *anypb.Any) error {
	ctx.Logger.Info("[worker handler for nodeC]")

	if ctx.IsMessageError() {
		res := &proto.Response{}
		res.Greeting = "Error during execution, please check previous logs"
		ctx.SendOutput(res)
	}

	if !attendedRequests[ctx.GetRequestID()] {
		ctx.SendAny(data)
	} else {
		delete(attendedRequests, ctx.GetRequestID())
	}

	return nil
}

func main() {
	handlers := map[string]kre.Handler{
		"nodeA": nodeA_handler,
		"nodeC": nodeC_handler,
	}

	kre.Start(handlerInit, defaultHandler, handlers)
}
