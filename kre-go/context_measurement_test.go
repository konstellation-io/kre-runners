package kre

import (
	"testing"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"

	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kre-runners/kre-go/config"
	"github.com/konstellation-io/kre-runners/kre-go/mocks"
	"github.com/konstellation-io/kre/libs/simplelogger"

	"github.com/stretchr/testify/suite"
)

type ContextMeasurementTestSuite struct {
	suite.Suite
	logger         *simplelogger.SimpleLogger
	mockController *gomock.Controller
	mockWriteAPI   *mocks.MockWriteAPI
	ctxMeasurement *contextMeasurement
}

func (suite *ContextMeasurementTestSuite) SetupSuite() {
	suite.logger = simplelogger.New(simplelogger.LevelInfo)
	suite.mockController = gomock.NewController(suite.T())
	suite.mockWriteAPI = mocks.NewMockWriteAPI(suite.mockController)

	suite.ctxMeasurement = &contextMeasurement{
		config.Config{
			Version:      "test_version",
			WorkflowName: "test_workflow",
			NodeName:     "test_node",
		},
		suite.logger,
		suite.mockWriteAPI,
	}
}

func TestContextMeasurementTestSuite(t *testing.T) {
	suite.Run(t, new(ContextMeasurementTestSuite))
}

func (suite *ContextMeasurementTestSuite) TestMeasurementSave() {
	measurement := "test_measurement"
	fields := map[string]interface{}{"field": "test"}
	tags := map[string]string{"tag": "test"}

	// do not use the monkey library aside from testing environment
	// here we need to patch through the time.Now() function so both timestamps are the same
	mockNow := time.Now()
	patch := monkey.Patch(time.Now, func() time.Time { return mockNow })
	defer patch.Unpatch()

	// make our own influx point, the one we are expecting will be written by the save function
	testPoint := influxdb2.NewPointWithMeasurement(measurement)
	testPoint.AddField("field", "test")
	testPoint.AddTag("tag", "test")
	testPoint.AddTag("version", suite.ctxMeasurement.cfg.Version)
	testPoint.AddTag("workflow", suite.ctxMeasurement.cfg.WorkflowName)
	testPoint.AddTag("node", suite.ctxMeasurement.cfg.NodeName)
	testPoint.SetTime(time.Now())

	suite.mockWriteAPI.EXPECT().WritePoint(testPoint).Return()
	suite.mockWriteAPI.EXPECT().Flush().Return()

	suite.ctxMeasurement.Save(measurement, fields, tags)
}
