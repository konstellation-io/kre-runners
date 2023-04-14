package main

import (
	"fmt"
	"math/rand"

	"main/proto"

	protobuf "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/konstellation-io/kre-runners/kre-go/v4"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

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

	if req.Testing.TestStores {
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

	testingResults := &proto.TestingResults{
		TestStoresSuccess: true,
	}

	//res.Greeting = randStringBytes(512 * 1024)
	res.Greeting = req.Greeting
	res.TestingResults = testingResults
	err = ctx.SendOutput(res)
	if err != nil {
		ctx.Logger.Errorf("Error sending output: %s", err)
		return err
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

	if string(objRes) != objValue {
		return fmt.Errorf("the value from the object store is not the same as the saved one")
	}

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

	if kvRes != objValue {
		return fmt.Errorf("the value from the kvStore is not the same as the saved one")
	}

	ctx.Configuration.Delete(kvKey)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	kre.Start(handlerInit, handler)
}
