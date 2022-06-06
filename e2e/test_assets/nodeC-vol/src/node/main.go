package main

import (
	"fmt"

	localProto "main/proto"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/konstellation-io/kre-runners/kre-go"
	proto2 "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func handlerInit(ctx *kre.HandlerContext) {
	ctx.Logger.Info("[worker init]")
}

func handler(ctx *kre.HandlerContext, data *any.Any) (proto.Message, error) {
	ctx.Logger.Info("[worker handler]")

	req := &localProto.NodeCRequest{}
	res := &localProto.Response{}

	ctx.Logger.Info(data.String())

	err := anypb.UnmarshalTo(data, req, proto2.UnmarshalOptions{})
	if err != nil {
		return res, fmt.Errorf("invalid request: %s", err)
	}

	result := req.Greeting + " from nodeC"

	ctx.Logger.Info(fmt.Sprintf("result -> %s", result))

	res.Greeting = result

	return res, nil
}

func main() {
	kre.Start(handlerInit, handler)
}
