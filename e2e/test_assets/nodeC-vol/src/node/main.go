package main

import (
	"fmt"

	"main/proto"

	protobuf "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/konstellation-io/kre-runners/kre-go/v4"
)

func handlerInit(ctx *kre.HandlerContext) {
	ctx.Logger.Info("[worker init]")
}

func handler(ctx *kre.HandlerContext, data *anypb.Any) error {
	ctx.Logger.Info("[worker handler]")

	req := &proto.NodeCRequest{}
	res := &proto.Response{}

	ctx.Logger.Info(data.String())

	err := anypb.UnmarshalTo(data, req, protobuf.UnmarshalOptions{})
	if err != nil {
		return fmt.Errorf("invalid request: %s", err)
	}

	result := req.Greeting + " and nodeC!"

	ctx.Logger.Info(fmt.Sprintf("result -> %s", result))

	res.Greeting = result

	ctx.SendOutput(res)

	testObjectStore(ctx)

	testKVStore(ctx)

	return nil
}

func testObjectStore(ctx *kre.HandlerContext) error {
	objKey := "42"
	objValue := "Test Node Object Store Config"

	err := ctx.StoreObject(objKey, []byte(objValue))
	if err != nil {
		return fmt.Errorf("error storing object: %w", err)
	}

	objRes, err := ctx.GetObject(objKey)
	if err != nil {
		return fmt.Errorf("error getting object: %w", err)
	}

	ctx.Logger.Info(string(objRes))

	ctx.DeleteObject(objKey)
	if err != nil {
		return fmt.Errorf("error deleting object: %w", err)
	}

	_, err = ctx.GetObject(objKey)
	if err != nil {
		ctx.Logger.Errorf("deleted key not found: %s", err)
	}

	return nil
}

func testKVStore(ctx *kre.HandlerContext) error {
	kvKey := "42"
	objValue := "Test Node KV Store Config"

	err := ctx.SetConfig(kvKey, objValue)
	if err != nil {
		return fmt.Errorf("error setting config: %w", err)
	}

	kvRes, err := ctx.GetConfig(kvKey, kre.ScopeNode)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	ctx.Logger.Info(kvRes)

	ctx.DeleteConfig(kvKey)
	if err != nil {
		return fmt.Errorf("error deleting config: %w", err)
	}

	_, err = ctx.GetConfig(kvKey, kre.ScopeNode)
	if err != nil {
		ctx.Logger.Errorf("deleted key not found: %s", err)
	}

	return nil
}

func main() {
	kre.Start(handlerInit, handler)
}
