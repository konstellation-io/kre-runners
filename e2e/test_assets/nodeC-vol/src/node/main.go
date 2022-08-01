package main

import (
	"fmt"

	"main/proto"

	protobuf "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/konstellation-io/kre-runners/kre-go"
)

func handlerInit(ctx *kre.HandlerContext) {
	ctx.Logger.Info("[worker init]")
}

func handler(ctx *kre.HandlerContext, data *anypb.Any) (protobuf.Message, error) {
	ctx.Logger.Info("[worker handler]")

	req := &proto.NodeCRequest{}
	res := &proto.Response{}

	ctx.Logger.Info(data.String())

	err := anypb.UnmarshalTo(data, req, protobuf.UnmarshalOptions{})
	if err != nil {
		return res, fmt.Errorf("invalid request: %s", err)
	}

	result := req.Greeting + " and nodeC!"

	ctx.Logger.Info(fmt.Sprintf("result -> %s", result))

	res.Greeting = result

	return res, nil
}

func main() {
	kre.Start(handlerInit, handler)
}
