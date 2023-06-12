package main

import (
	"fmt"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/runner/trigger"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/context"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/runner"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/anypb"
)

var wg sync.WaitGroup

func main() {
	fmt.Println("Initializing NATS config...")
	fmt.Println()
	go publishMessageToNats()

	time.Sleep(10 * time.Second)

	fmt.Println("Initializing Trigger Runner...")
	fmt.Println()

	go runner.
		NewRunner().
		TriggerRunner().
		WithInitializer(func(ctx context.KaiContext) {
			err := ctx.ObjectStore.Save("test", []byte("testValue"))
			if err != nil {
				ctx.Logger.Error(err, "Error saving object")
			}

			err = ctx.CentralizedConfig.SetConfig("test", "testConfigValue")
			if err != nil {
				ctx.Logger.Error(err, "Error setting config")
			}

			ctx.Logger.Info("Metadata",
				"process", ctx.Metadata.GetProcess(),
				"product", ctx.Metadata.GetProduct(),
				"workflow", ctx.Metadata.GetWorkflow(),
				"version", ctx.Metadata.GetVersion(),
				"kv_product", ctx.Metadata.GetKeyValueStoreProductName(),
				"kv_workflow", ctx.Metadata.GetKeyValueStoreWorkflowName(),
				"kv_process", ctx.Metadata.GetKeyValueStoreProcessName(),
				"object_store", ctx.Metadata.GetObjectStoreName(),
			)

			ctx.Logger.Info("PathUtils",
				"getBasePath", ctx.PathUtils.GetBasePath(),
				"composeBasePath", ctx.PathUtils.ComposeBasePath("test"))
		}).
		WithRunner(func(ctx context.KaiContext, handlers map[string]trigger.ResponseHandler) {
			http.HandleFunc("/hello", responseHandler(ctx, handlers))
			ctx.Logger.Info("Starting http server", "port", 8080)
			ctx.Logger.Error(http.ListenAndServe(":8080", nil), "Error serving http")
		}).
		WithFinalizer(func(ctx context.KaiContext) {
			ctx.Logger.Info("Finalizer")
			wg.Done()
		}).
		Run()

	time.Sleep(5 * time.Second)
	fmt.Println("Initializing Exit Runner...")
	fmt.Println()

	/*go runner.
		NewRunner().
	]TaskRunner().
		WithInitializer(func(ctx context.KaiContext) {
			value, err := ctx.CentralizedConfig.GetConfig("test")
			if err != nil {
				ctx.Logger.Error(err, "Error getting config")
				return
			}
			ctx.Logger.Info("Config value retrieved!", "value", value)

			obj, err := ctx.ObjectStore.Get("test")
			if err != nil {
				ctx.Logger.Error(err, "Error getting Obj Store values")
				return
			}
			ctx.Logger.Info("ObjectStore value retrieved!", "object", string(obj))
		}).
		WithPreprocessor(func(ctx context.KaiContext, response *anypb.Any) error {
			ctx.Logger.Info("Preprocessing event")
			return nil
		}).
		WithHandler(func(ctx context.KaiContext, response *anypb.Any) error {
			ctx.Logger.Info("Response Handler, message received",
				"response", string(response.GetValue()))
			return nil
		}).
		WithPostprocessor(func(ctx context.KaiContext, response *anypb.Any) error {
			ctx.Logger.Info("Postprocessor event")
			return nil
		}).
		WithFinalizer(func(ctx context.KaiContext) {
			ctx.Logger.Info("Finalizer")
			wg.Done()
		}).
		Run()

	fmt.Println()
	fmt.Println()*/

	go runner.
		NewRunner().
		ExitRunner().
		WithInitializer(func(ctx context.KaiContext) {
			value, err := ctx.CentralizedConfig.GetConfig("test")
			if err != nil {
				ctx.Logger.Error(err, "Error getting config")
				return
			}
			ctx.Logger.Info("Config value retrieved!", "value", value)

			obj, err := ctx.ObjectStore.Get("test")
			if err != nil {
				ctx.Logger.Error(err, "Error getting Obj Store values")
				return
			}
			ctx.Logger.Info("ObjectStore value retrieved!", "object", string(obj))
		}).
		WithPreprocessor(func(ctx context.KaiContext, response *anypb.Any) error {
			ctx.Logger.Info("Preprocessor")
			return nil
		}).
		WithHandler(func(ctx context.KaiContext, response *anypb.Any) error {
			ctx.Logger.Info("Handler")
			// ctx.Messaging.SendAny(response)
			return nil
		}).
		WithPostprocessor(func(ctx context.KaiContext, response *anypb.Any) error {
			ctx.Logger.Info("Postprocessor")
			return nil
		}).
		WithFinalizer(func(ctx context.KaiContext) {
			ctx.Logger.Info("Finalizer")
			wg.Done()
		}).
		Run()

	wg.Add(2)
	wg.Wait()
}

func responseHandler(ctx context.KaiContext, handlers map[string]trigger.ResponseHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		nameRequest := r.URL.Query().Get("name")

		stringb := []byte(fmt.Sprintf("Hello %s!", nameRequest))

		anyValue := &anypb.Any{
			TypeUrl: "type.googleapis.com/google.protobuf.StringValue",
			Value:   stringb,
		}

		reqID := uuid.NewString()

		ctx.Logger.Info("Sending message to nats",
			"execution ID", reqID, "message", anyValue)

		err := ctx.Messaging.SendOutputWithRequestID(anyValue, reqID)
		if err != nil {
			ctx.Logger.Error(err, "Error sending message to nats")
			return
		}

		var wg2 sync.WaitGroup
		wg2.Add(1)

		handlers[reqID] = func(ctx context.KaiContext, response *anypb.Any) error {
			defer wg2.Done()
			_, err2 := w.Write([]byte(fmt.Sprintf("Hello, World! %s", response.GetValue())))
			if err2 != nil {
				return err2
			}

			return nil
		}

		wg2.Wait()
	}
}

func publishMessageToNats() {
	nc, err := nats.Connect(viper.GetString("NATS_URL"))
	if err != nil {
		panic(err)
	}

	// Use the JetStream context to produce and consumer messages
	// that have been persisted.
	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}

	js.AddStream(&nats.StreamConfig{
		Name: "test-stream",
		Subjects: []string{
			"test-input",
			"test-output",
		},
		Retention: nats.InterestPolicy,
	})

	js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket:  "test1",
		Storage: nats.MemoryStorage,
	})
	js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket:  "test2",
		Storage: nats.MemoryStorage,
	})
	js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket:  "test3",
		Storage: nats.MemoryStorage,
	})

	js.CreateObjectStore(&nats.ObjectStoreConfig{
		Bucket: "test",
	})
}
