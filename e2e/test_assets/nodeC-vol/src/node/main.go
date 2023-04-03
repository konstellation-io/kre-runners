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
		err = testObjectStore(ctx)
		if err != nil {
			return fmt.Errorf("error during test object store: %w", err)
		}

		ctx.Logger.Info("testing configuration")
		err = testKVStore(ctx)
		if err != nil {
			return fmt.Errorf("error during test key value store: %w", err)
		}
	}
	return nil
}

func testObjectStore(ctx *kre.HandlerContext) error {
	objKey := "test_object_key"
	objValue := fmt.Sprintf("%d", rand.Intn(100000))

	err := ctx.ObjectStore.Save(objKey, []byte(objValue))
	if err != nil {
		return err
	}

	objRes, err := ctx.ObjectStore.Get(objKey)
	if err != nil {
		return err
	}

	ctx.Logger.Infof("the value from the object store is: %s", string(objRes))

	ctx.ObjectStore.Delete(objKey)
	if err != nil {
		return err
	}

	return nil
}

func testKVStore(ctx *kre.HandlerContext) error {
	kvKey := "test_kv_key"
	objValue := fmt.Sprintf("Test Node Value: %d", rand.Intn(100000))

	err := ctx.Configuration.Set(kvKey, objValue)
	if err != nil {
		return err
	}

	kvRes, err := ctx.Configuration.Get(kvKey, kre.NodeScope)
	if err != nil {
		return err
	}

	ctx.Logger.Infof("the config from the kvStore is: %s", kvRes)

	ctx.Configuration.Delete(kvKey)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	kre.Start(handlerInit, handler)
}
