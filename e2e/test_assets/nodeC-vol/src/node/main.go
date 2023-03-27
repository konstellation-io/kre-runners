package main

import (
	"fmt"
	"math/rand"

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

	if randomInt := rand.Intn(100); randomInt < 5 {
		ctx.Logger.Info("testing object store")
		testObjectStore(ctx)

		ctx.Logger.Info("testing configuration")
		testKVStore(ctx)
	}
	return nil
}

func testObjectStore(ctx *kre.HandlerContext) error {
	objKey := "42"
	objValue := "Test Node Object Store Config"

	err := ctx.ObjectStore.Save(objKey, []byte(objValue))
	if err != nil {
		return fmt.Errorf("error storing object: %w", err)
	}

	objRes, err := ctx.ObjectStore.Get(objKey)
	if err != nil {
		return fmt.Errorf("error getting object: %w", err)
	}

	ctx.Logger.Infof("the value from the object store is: %s", string(objRes))

	ctx.ObjectStore.Delete(objKey)
	if err != nil {
		return fmt.Errorf("error deleting object: %w", err)
	}

	_, err = ctx.ObjectStore.Get(objKey)
	if err != nil {
		ctx.Logger.Errorf("deleted key not found: %s", err)
	}

	return nil
}

func testKVStore(ctx *kre.HandlerContext) error {
	kvKey := "42"
	objValue := "Test Node KV Store Config"

	err := ctx.Configuration.Set(kvKey, objValue)
	if err != nil {
		return fmt.Errorf("error setting config: %w", err)
	}

	kvRes, err := ctx.Configuration.Get(kvKey, kre.ScopeNode)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	ctx.Logger.Infof("the config from the kvStore is %s", kvRes)

	ctx.Configuration.Delete(kvKey)
	if err != nil {
		return fmt.Errorf("error deleting config: %w", err)
	}

	_, err = ctx.Configuration.Get(kvKey, kre.ScopeNode)
	if err != nil {
		ctx.Logger.Errorf("deleted key not found: %s", err)
	}

	return nil
}

func main() {
	kre.Start(handlerInit, handler)
}
