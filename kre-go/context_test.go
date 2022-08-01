package kre

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kre/libs/simplelogger"
	testserver "github.com/nats-io/nats-server/v2/test"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/konstellation-io/kre-runners/kre-go/config"
	"github.com/konstellation-io/kre-runners/kre-go/mocks"
)

type TestPrediction struct {
	Time       time.Time
	Prediction string
	TicketID   string
	Asset      string
}

func setEnvVars(t *testing.T, envVars map[string]string) {
	for name, value := range envVars {
		err := os.Setenv(name, value)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestHandlerContext_GetData(t *testing.T) {
	logger := simplelogger.New(simplelogger.LevelDebug)
	cfg := config.Config{}

	testPort := 8331
	opts := testserver.DefaultTestOptions
	opts.Port = testPort
	s := testserver.RunServer(&opts)
	defer s.Shutdown()

	nc, err := nats.Connect(fmt.Sprintf("nats://127.0.0.1:%d", testPort))
	if err != nil {
		t.Fatal(err)
	}
	defer nc.Close()

	msgCh := make(chan *nats.Msg, 64)

	sub, err := nc.ChanSubscribe("mongo_writer", msgCh)
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mongoM := mocks.NewMockManager(ctrl)

	ctx := NewHandlerContext(cfg, nc, mongoM, logger, nil, nil)

	q := QueryData{
		"TicketID": "1234",
		"Asset":    "A5678",
	}

	var results []*TestPrediction

	savedData := TestPrediction{
		Time:       time.Now(),
		Prediction: "Repair Complete",
		TicketID:   "1234",
		Asset:      "A5678",
	}

	criteria := bson.M{
		"TicketID": "1234",
		"Asset":    "A5678",
	}

	mongoM.EXPECT().
		Find(gomock.Any(), "test_predictions", criteria, results).
		Return(nil).
		Do(func(ctx context.Context, colName string, filter interface{}, _ interface{}) {
			results = append(results, &TestPrediction{
				Time:       time.Now(),
				Prediction: "Repair Complete",
				TicketID:   "1234",
				Asset:      "A5678",
			})
		})

	err = ctx.DB.Find("test_predictions", q, results)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) == 0 {
		t.Fatalf("no result returned")
	}
	result := results[0]
	if result.Prediction != savedData.Prediction {
		t.Fatalf("saved data is wrong: %s != %s", result.Prediction, savedData.Prediction)
	}
}

func TestHandlerContext_SaveData(t *testing.T) {
	const mongoWriterSubject = "mongo_writer"

	logger := simplelogger.New(simplelogger.LevelDebug)

	cfg := config.Config{
		NATS: config.ConfigNATS{
			MongoWriterSubject: mongoWriterSubject,
		},
		MongoDB: config.MongoDB{
			Address: "mongodb://localhost:27017",
		},
	}

	testPort := 8331
	opts := testserver.DefaultTestOptions
	opts.Port = testPort
	s := testserver.RunServer(&opts)
	defer s.Shutdown()

	nc, err := nats.Connect(fmt.Sprintf("nats://127.0.0.1:%d", testPort))
	if err != nil {
		t.Fatal(err)
	}
	defer nc.Close()

	msgCh := make(chan *nats.Msg, 64)

	sub, err := nc.ChanSubscribe(mongoWriterSubject, msgCh)
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	ctrl := gomock.NewController(t)
	mongoM := mocks.NewMockManager(ctrl)

	ctx := NewHandlerContext(cfg, nc, mongoM, logger, nil, nil)
	sentMsg := TestPrediction{
		Time:       time.Now(),
		Prediction: "Tested",
		TicketID:   "1234",
		Asset:      "A12345C",
	}

	c, cancel := context.WithCancel(context.Background())

	var goFuncErr error
	go func() {
		err := ctx.DB.Save("test_predictions", sentMsg)
		goFuncErr = err
		cancel()
	}()

	receivedMsg := SaveDataMsg{}

	msg := <-msgCh
	if goFuncErr != nil {
		t.Fatal(goFuncErr)
	}

	err = json.Unmarshal(msg.Data, &receivedMsg)
	if err != nil {
		t.Fatalf("Error parsing data msg: %s", err)
	}

	err = msg.Respond([]byte("{ Success: true }"))
	if err != nil {
		t.Fatalf("Error replaying to the data msg: %s", err)
	}

	<-c.Done()

	receivedDoc := receivedMsg.Doc.(map[string]interface{})

	if receivedDoc["TicketID"] != sentMsg.TicketID {
		t.Fatalf("TicketID value is wrong %s != %s", receivedDoc["TicketID"], sentMsg.TicketID)
	}
	if receivedDoc["Asset"] != sentMsg.Asset {
		t.Fatalf("Asset value is wrong %s != %s", receivedDoc["Asset"], sentMsg.Asset)
	}
}
