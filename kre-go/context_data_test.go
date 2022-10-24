package kre

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nats-io/nats-server/v2/server"
	testserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/suite"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/konstellation-io/kre-runners/kre-go/config"
	"github.com/konstellation-io/kre-runners/kre-go/mocks"
	"github.com/konstellation-io/kre/libs/simplelogger"
)

type TestData struct {
	Time     time.Time
	Result   string
	TicketID string
	Asset    string
}

type ContextDataTestSuite struct {
	suite.Suite
	tServer *server.Server
	nc      *nats.Conn
	mongoM  *mocks.MockManager
	ctrl    *gomock.Controller
	ctxData *contextData
}

const mongoWriterSubject = "mongo_writer"

func TestContextDataTestSuite(t *testing.T) {
	suite.Run(t, new(ContextDataTestSuite))
}

func (suite *ContextDataTestSuite) SetupSuite() {
	logger := simplelogger.New(simplelogger.LevelInfo)
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
	suite.tServer = testserver.RunServer(&opts)

	var err error
	suite.nc, err = nats.Connect(fmt.Sprintf("nats://127.0.0.1:%d", testPort))
	suite.Require().NoError(err)

	suite.ctrl = gomock.NewController(suite.T())
	suite.mongoM = mocks.NewMockManager(suite.ctrl)

	suite.ctxData = NewContextData(cfg, suite.nc, suite.mongoM, logger)
}

func (suite *ContextDataTestSuite) TearDownSuite() {
	suite.ctrl.Finish()
	suite.nc.Close()
	suite.tServer.Shutdown()
}

func (suite *ContextDataTestSuite) TestContextDataFind() {
	qData := QueryData{
		"TicketID": "1234",
		"Asset":    "A5678",
	}

	var results []*TestData

	savedData := TestData{
		Time:     time.Now(),
		Result:   "Repair Complete",
		TicketID: "1234",
		Asset:    "A5678",
	}

	criteria := bson.M{
		"TicketID": "1234",
		"Asset":    "A5678",
	}

	suite.mongoM.EXPECT().
		Find(gomock.Any(), "test_data", criteria, results).
		Return(nil).
		Do(func(ctx context.Context, colName string, filter interface{}, _ interface{}) {
			results = append(results, &TestData{
				Time:     time.Now(),
				Result:   "Repair Complete",
				TicketID: "1234",
				Asset:    "A5678",
			})
		})

	err := suite.ctxData.Find("test_data", qData, results)
	suite.Require().NoError(err)

	suite.Require().Len(results, 1)
	result := results[0]
	suite.Equal(savedData.Result, result.Result)
}

func (suite *ContextDataTestSuite) TestContextDataSave() {
	msgCh := make(chan *nats.Msg, 64)
	sub, err := suite.nc.ChanSubscribe(mongoWriterSubject, msgCh)
	suite.Require().NoError(err)
	defer sub.Unsubscribe()

	sentMsg := TestData{
		Time:     time.Now(),
		Result:   "Tested",
		TicketID: "1234",
		Asset:    "A12345C",
	}

	c, cancel := context.WithCancel(context.Background())

	var goFuncErr error
	go func() {
		err := suite.ctxData.Save("test_data", sentMsg)
		goFuncErr = err
		cancel()
	}()

	receivedMsg := SaveDataMsg{}

	msg := <-msgCh
	suite.Require().NoError(goFuncErr)

	err = json.Unmarshal(msg.Data, &receivedMsg)
	suite.Require().NoError(err)

	err = msg.Respond([]byte("{ Success: true }"))
	suite.Require().NoError(err)

	<-c.Done()

	receivedDoc := receivedMsg.Doc.(map[string]interface{})

	suite.Assert().Equal(sentMsg.TicketID, receivedDoc["TicketID"])
	suite.Assert().Equal(sentMsg.Asset, receivedDoc["Asset"])
}
