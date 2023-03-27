package kre

import (
	"fmt"
	"testing"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/suite"

	testserver "github.com/nats-io/nats-server/v2/test"

	"github.com/konstellation-io/kre-runners/kre-go/v4/config"
	"github.com/konstellation-io/kre/libs/simplelogger"
)

type ContextConfigurationTestSuite struct {
	suite.Suite
	tServer          *server.Server
	logger           *simplelogger.SimpleLogger
	cfg              config.Config
	nc               *nats.Conn
	js               nats.JetStreamContext
	ctxConfiguration *contextConfiguration
}

func TestContextConfigurationTestSuite(t *testing.T) {
	suite.Run(t, new(ContextConfigurationTestSuite))
}

// SetupSuite will create a logger, setup a config, run a NATS mocked server,
// then connect a NATS client to it to get a JetStream context.
//
// These will be used to create a context configuration object.
func (suite *ContextConfigurationTestSuite) SetupSuite() {
	suite.logger = simplelogger.New(simplelogger.LevelInfo)
	suite.cfg = config.Config{
		NATS: config.ConfigNATS{
			KeyValueStoreProjectName:  "kv_project",
			KeyValueStoreWorkflowName: "kv_workflow",
			KeyValueStoreNodeName:     "kv_node",
		},
	}

	testPort := 8331
	opts := testserver.DefaultTestOptions
	opts.Port = testPort
	opts.JetStream = true
	suite.tServer = testserver.RunServer(&opts)

	var err error
	suite.nc, err = nats.Connect(fmt.Sprintf("nats://127.0.0.1:%d", testPort))
	suite.Require().NoError(err)

	suite.js, err = suite.nc.JetStream()
	suite.Require().NoError(err)

	suite.createKVStores(suite.cfg)
}

// TearDownSuite will close the mock controller, close the NATS connection and shutdown the mocked server
func (suite *ContextConfigurationTestSuite) TearDownSuite() {
	suite.nc.Close()
	suite.tServer.Shutdown()
}

// SetupTest will run before each test
func (suite *ContextConfigurationTestSuite) SetupTest() {
	var err error
	suite.ctxConfiguration, err = NewContextConfiguration(suite.cfg, suite.logger, suite.js)
	suite.Require().NoError(err)
}

func (suite *ContextConfigurationTestSuite) createKVStores(cfg config.Config) {
	allKVStores := []string{
		cfg.NATS.KeyValueStoreProjectName,
		cfg.NATS.KeyValueStoreWorkflowName,
		cfg.NATS.KeyValueStoreNodeName,
	}

	for _, kvStore := range allKVStores {
		cfg := nats.KeyValueConfig{
			Bucket:  kvStore,
			Storage: nats.FileStorage,
		}
		_, err := suite.js.CreateKeyValue(&cfg)
		suite.Require().NoError(err)
	}
}

func (suite *ContextConfigurationTestSuite) getValueFromKVStore(bucket, key string) string {
	KVStore, err := suite.js.KeyValue(bucket)
	suite.Require().NoError(err)

	entry, err := KVStore.Get(key)
	suite.Require().NoError(err)

	return string(entry.Value())
}

func (suite *ContextConfigurationTestSuite) TestNewContextConfiguration() {
	suite.Require().NotNil(suite.ctxConfiguration)
	suite.Require().NotNil(suite.ctxConfiguration.kvStoresMap)
	suite.Require().Len(suite.ctxConfiguration.kvStoresMap, 3)
}

func (suite *ContextConfigurationTestSuite) TestSetProjectConfig() {
	key := "key"
	value := "value"

	// to project bucket
	err := suite.ctxConfiguration.Set(key, value, ScopeProject)
	suite.Require().NoError(err)
	savedValue := suite.getValueFromKVStore(suite.cfg.NATS.KeyValueStoreProjectName, key)
	suite.Assert().Equal(value, savedValue)

	// to workflow bucket
	err = suite.ctxConfiguration.Set(key, value, ScopeWorkflow)
	suite.Require().NoError(err)
	savedValue = suite.getValueFromKVStore(suite.cfg.NATS.KeyValueStoreWorkflowName, key)
	suite.Assert().Equal(value, savedValue)

	// to node bucket
	value = "value2"
	err = suite.ctxConfiguration.Set(key, value, ScopeNode)
	suite.Require().NoError(err)
	savedValue = suite.getValueFromKVStore(suite.cfg.NATS.KeyValueStoreNodeName, key)
	suite.Assert().Equal(value, savedValue)

	// to deafult bucket (node)
	err = suite.ctxConfiguration.Set(key, value)
	suite.Require().NoError(err)
	savedValue = suite.getValueFromKVStore(suite.cfg.NATS.KeyValueStoreNodeName, key)
	suite.Assert().Equal(value, savedValue)
}

func (suite *ContextConfigurationTestSuite) TestGetProjectConfig() {
	key := "key"
	value := "value"

	// from project bucket
	err := suite.ctxConfiguration.Set(key, value, ScopeProject)
	suite.Require().NoError(err)
	savedValue, err := suite.ctxConfiguration.Get(key, ScopeProject)
	suite.Require().NoError(err)
	suite.Assert().Equal(value, savedValue)

	// from workflow bucket
	err = suite.ctxConfiguration.Set(key, value, ScopeWorkflow)
	suite.Require().NoError(err)
	savedValue, err = suite.ctxConfiguration.Get(key, ScopeWorkflow)
	suite.Require().NoError(err)
	suite.Assert().Equal(value, savedValue)

	// from node bucket
	err = suite.ctxConfiguration.Set(key, value, ScopeNode)
	suite.Require().NoError(err)
	savedValue, err = suite.ctxConfiguration.Get(key, ScopeNode)
	suite.Require().NoError(err)
	suite.Assert().Equal(value, savedValue)

	// from default bucket (node)
	value = "value2"
	err = suite.ctxConfiguration.Set(key, value)
	suite.Require().NoError(err)
	savedValue, err = suite.ctxConfiguration.Get(key)
	suite.Require().NoError(err)
	suite.Assert().Equal(value, savedValue)
}

func (suite *ContextConfigurationTestSuite) TestDeleteProjectConfig() {
	key := "key"
	value := "value"

	// from project bucket
	err := suite.ctxConfiguration.Set(key, value, ScopeProject)
	suite.Require().NoError(err)
	err = suite.ctxConfiguration.Delete(key, ScopeProject)
	suite.Require().NoError(err)
	_, err = suite.ctxConfiguration.Get(key, ScopeProject)
	suite.Require().Error(err)

	// from workflow bucket
	err = suite.ctxConfiguration.Set(key, value, ScopeWorkflow)
	suite.Require().NoError(err)
	err = suite.ctxConfiguration.Delete(key, ScopeWorkflow)
	suite.Require().NoError(err)
	_, err = suite.ctxConfiguration.Get(key, ScopeWorkflow)
	suite.Require().Error(err)

	// from node bucket
	err = suite.ctxConfiguration.Set(key, value, ScopeNode)
	suite.Require().NoError(err)
	err = suite.ctxConfiguration.Delete(key, ScopeNode)
	suite.Require().NoError(err)
	_, err = suite.ctxConfiguration.Get(key, ScopeNode)
	suite.Require().Error(err)

	// from default bucket (node)
	err = suite.ctxConfiguration.Set(key, value)
	suite.Require().NoError(err)
	err = suite.ctxConfiguration.Delete(key)
	suite.Require().NoError(err)
	_, err = suite.ctxConfiguration.Get(key)
	suite.Require().Error(err)
}
