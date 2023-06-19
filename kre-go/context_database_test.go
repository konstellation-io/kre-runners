//go:build integration

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

	"github.com/konstellation-io/kre-runners/kre-go/v4/config"
	"github.com/konstellation-io/kre-runners/kre-go/v4/mocks"
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
	ctxData ContextDatabase
}

const mongoWriterSubject = "mongo_writer"

func TestContextDataTestSuite(t *testing.T) {
	suite.Run(t, new(ContextDataTestSuite))
}

// SetupSuite will create a logger, setup a config, run a NATS mocked server, then connect a
// NATS client to it, also create a mock controller and a mock Mongo manager.
// These will be used to create a context data object.
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

	suite.ctxData = NewContextDatabase(cfg, suite.nc, suite.mongoM, logger)
}

// TearDownSuite will close the mock controller, close the NATS connection and shutdown the mocked server
func (suite *ContextDataTestSuite) TearDownSuite() {
	suite.ctrl.Finish()
	suite.nc.Close()
	suite.tServer.Shutdown()
}

func (suite *ContextDataTestSuite) TestContextDataFind() {
	// GIVEN  a query
	qData := QueryData{
		"TicketID": "1234",
		"Asset":    "A5678",
	}

	// GIVEN a collection to where load the query results
	var results []*TestData

	// GIVEN the data we expect to collect
	savedData := TestData{
		Time:     time.Now(),
		Result:   "Repair Complete",
		TicketID: "1234",
		Asset:    "A5678",
	}

	// GIVEN the criteria search we expect will be used
	criteria := bson.M{
		"TicketID": "1234",
		"Asset":    "A5678",
	}

	// WHEN the find method is called
	// THEN append the expected result to the results collection as if the function returned this data
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

	// THEN returned data must be of length 1 and equal to the data we expect to be obtained
	suite.Require().Len(results, 1)
	result := results[0]
	suite.Equal(savedData.Result, result.Result)
}

func (suite *ContextDataTestSuite) TestContextDataSave() {
	// GIVEN  channel subscription to the mongo writter
	msgCh := make(chan *nats.Msg, 64)
	sub, err := suite.nc.ChanSubscribe(mongoWriterSubject, msgCh)
	suite.Require().NoError(err)
	defer sub.Unsubscribe()

	// GIVEN some test data
	sentMsg := TestData{
		Time:     time.Now(),
		Result:   "Tested",
		TicketID: "1234",
		Asset:    "A12345C",
	}

	c, cancel := context.WithCancel(context.Background())

	// WHEN test information is saved
	var goFuncErr error
	go func() {
		err := suite.ctxData.Save("test_data", sentMsg)
		goFuncErr = err
		cancel()
	}()

	receivedMsg := SaveDataMsg{}

	// THEN expect a message written into the mongo writter topic with no error
	msg := <-msgCh
	suite.Require().NoError(goFuncErr)

	// THEN the message can be unmarshalled with no error
	err = json.Unmarshal(msg.Data, &receivedMsg)
	suite.Require().NoError(err)

	// THEN expect the message can be answered back with no error
	err = msg.Respond([]byte("{ Success: true }"))
	suite.Require().NoError(err)

	<-c.Done()

	// THEN the unmarshalled message equals to the test info we are expecting
	receivedDoc := receivedMsg.Doc.(map[string]interface{})

	suite.Equal(sentMsg.TicketID, receivedDoc["TicketID"])
	suite.Equal(sentMsg.Asset, receivedDoc["Asset"])
}
