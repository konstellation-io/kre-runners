package kre

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/konstellation-io/kre/libs/simplelogger"

	"github.com/nats-io/nats.go"

	"github.com/konstellation-io/kre-runners/kre-go/config"
)

func handler(ctx *HandlerContext, data *any.Any) error {
	ctx.Logger.Info("[worker handler]")

	input := &TestInput{}
	err := ptypes.UnmarshalAny(data, input)
	if err != nil {
		return err
	}

	greetingText := fmt.Sprintf("%s %s!", ctx.Get("greeting"), input.Name)
	ctx.Logger.Info(greetingText)

	out := &TestOutput{}
	out.Greeting = greetingText
	ctx.SendOutput(out)
	return nil
}

func setEnvVars(t *testing.T, envVars map[string]string) {
	for name, value := range envVars {
		err := os.Setenv(name, value)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestStart(t *testing.T) {
	const inputSubject = "test-subject-input"
	setEnvVars(t, map[string]string{
		"KRT_WORKFLOW_NAME":     "workflowTest",
		"KRT_VERSION":           "testVersion1",
		"KRT_VERSION_ID":        "version.12345",
		"KRT_NODE_NAME":         "nodeA",
		"KRT_BASE_PATH":         "./test",
		"KRT_NATS_SERVER":       "localhost:4222",
		"KRT_NATS_INPUT":        inputSubject,
		"KRT_NATS_OUTPUT":       "",
		"KRT_NATS_MONGO_WRITER": "mongo_writer",
		"KRT_MONGO_URI":         "mongodb://localhost:27017",
		"KRT_INFLUX_URI":        "influxdb-uri",
	})
	logger := simplelogger.New(simplelogger.LevelDebug)

	cfg := config.NewConfig(logger)
	nc, err := nats.Connect(cfg.NATS.Server)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	input := &TestInput{Name: "John"}
	inputData, err := ptypes.MarshalAny(input)
	if err != nil {
		t.Fatal(err)
	}

	kreNatsMsg := &KreNatsMessage{
		TrackingId: "msg.12345",
		Payload:    inputData,
		Tracking: []*KreNatsMessage_Tracking{
			{NodeName: "nodeTest", Start: time.Now().Format(ISO8601), End: time.Now().Format(ISO8601)},
		},
	}
	msg, err := proto.Marshal(kreNatsMsg)
	if err != nil {
		t.Fatal(err)
	}

	// Subscribe to mongo_writer subject and reply to save_metrics and save_data requests
	sub, err := nc.Subscribe("mongo_writer", func(msg *nats.Msg) {
		_ = nc.Publish(msg.Reply, []byte("{\"success\":true}"))
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	doneCh := make(chan struct{})

	handlerInit := func(ctx *HandlerContext) {
		ctx.Logger.Info("[worker init]")
		ctx.Set("greeting", "Hello")
		doneCh <- struct{}{}
	}
	go Start(handlerInit, handler)

	<-doneCh

	res, err := nc.Request(inputSubject, msg, 10*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	expectedRes := "Hello John!"
	result := &KreNatsMessage{}
	err = proto.Unmarshal(res.Data, result)
	if err != nil {
		t.Fatal(err)
	}

	resultData := &TestOutput{}
	err = ptypes.UnmarshalAny(result.Payload, resultData)
	if err != nil {
		t.Fatal(err)
	}

	if expectedRes != resultData.Greeting {
		t.Fatalf("The greeting '%s' must be '%s'", resultData.Greeting, expectedRes)
	}
}
