package kre

import (
	"testing"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"

	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kre-runners/kre-go/v4/config"
	"github.com/konstellation-io/kre-runners/kre-go/v4/mocks"
	"github.com/konstellation-io/kre/libs/simplelogger"

	"github.com/stretchr/testify/suite"
)

type ContextMeasurementTestSuite struct {
	suite.Suite
	mockWriteAPI   *mocks.MockWriteAPI
	ctxMeasurement *contextMeasurement
}

func TestContextMeasurementTestSuite(t *testing.T) {
	suite.Run(t, new(ContextMeasurementTestSuite))
}

// SetupSuite will create a mock controller and mock write api for influx.
// These will be also attached to a custom generated context measurement object
func (suite *ContextMeasurementTestSuite) SetupSuite() {
	logger := simplelogger.New(simplelogger.LevelInfo)
	mockController := gomock.NewController(suite.T())
	suite.mockWriteAPI = mocks.NewMockWriteAPI(mockController)

	suite.ctxMeasurement = &contextMeasurement{ // cannot use NewContextMeasurement as it initializes its own writeAPI
		config.Config{
			Version:      "test_version",
			WorkflowName: "test_workflow",
			NodeName:     "test_node",
		},
		logger,
		suite.mockWriteAPI,
	}
}

func (suite *ContextMeasurementTestSuite) TestContextMeasurementSave() {
	// GIVEN a new created metric
	measurement := "test_measurement"
	fields := map[string]interface{}{"field": "test"}
	tags := map[string]string{"tag": "test"}

	// do not use the monkey library aside from testing environment
	// here we need to patch through the time.Now() function so both timestamps are the same
	mockNow := time.Now()
	patch := monkey.Patch(time.Now, func() time.Time { return mockNow })
	defer patch.Unpatch()

	// GIVEN the metric we are expecting will be written by the save function
	testPoint := influxdb2.NewPointWithMeasurement(measurement)
	testPoint.AddField("field", "test")
	testPoint.AddTag("tag", "test")
	testPoint.AddTag("version", suite.ctxMeasurement.cfg.Version)
	testPoint.AddTag("workflow", suite.ctxMeasurement.cfg.WorkflowName)
	testPoint.AddTag("node", suite.ctxMeasurement.cfg.NodeName)
	testPoint.SetTime(time.Now())

	// WHEN the metric is saved
	// THEN the save method will be called once, with a measurement exact to our test point
	suite.mockWriteAPI.EXPECT().WritePoint(testPoint).Times(1).Return()
	suite.mockWriteAPI.EXPECT().Flush().Times(1).Return()

	suite.ctxMeasurement.Save(measurement, fields, tags)
}
