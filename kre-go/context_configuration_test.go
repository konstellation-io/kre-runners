//go:build integration

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

var (
	testKey     = "key"
	testValue   = "value"
	allKVStores = make([]string, 3)
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

func (suite *ContextConfigurationTestSuite) SetupSuite() {
	suite.logger = simplelogger.New(simplelogger.LevelInfo)
	suite.cfg = config.Config{
		NATS: config.ConfigNATS{
			KeyValueStoreProjectName:  "kv_project",
			KeyValueStoreWorkflowName: "kv_workflow",
			KeyValueStoreNodeName:     "kv_node",
		},
	}

	allKVStores = []string{
		suite.cfg.NATS.KeyValueStoreProjectName,
		suite.cfg.NATS.KeyValueStoreWorkflowName,
		suite.cfg.NATS.KeyValueStoreNodeName,
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
}

func (suite *ContextConfigurationTestSuite) TearDownSuite() {
	suite.nc.Close()
	suite.tServer.Shutdown()
}

func (suite *ContextConfigurationTestSuite) SetupTest() {
	var err error

	suite.createKVStores()

	suite.ctxConfiguration, err = NewContextConfiguration(suite.cfg, suite.logger, suite.js)
	suite.Require().NoError(err)
}

func (suite *ContextConfigurationTestSuite) TearDownTest() {
	suite.deleteKVStores()
}

func (suite *ContextConfigurationTestSuite) createKVStores() {
	for _, kvStore := range allKVStores {
		cfg := nats.KeyValueConfig{
			Bucket:  kvStore,
			Storage: nats.FileStorage,
		}
		_, err := suite.js.CreateKeyValue(&cfg)
		suite.Require().NoError(err)
	}
}

func (suite *ContextConfigurationTestSuite) deleteKVStores() {
	for _, kvStore := range allKVStores {
		err := suite.js.DeleteKeyValue(kvStore)
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

func (suite *ContextConfigurationTestSuite) TestSetConfigProjectScope() {
	err := suite.ctxConfiguration.Set(testKey, testValue, ProjectScope)
	suite.Require().NoError(err)
	savedValue := suite.getValueFromKVStore(suite.cfg.NATS.KeyValueStoreProjectName, testKey)
	suite.Assert().Equal(testValue, savedValue)
}

func (suite *ContextConfigurationTestSuite) TestSetConfigWorkflowScope() {
	err := suite.ctxConfiguration.Set(testKey, testValue, WorkflowScope)
	suite.Require().NoError(err)
	savedValue := suite.getValueFromKVStore(suite.cfg.NATS.KeyValueStoreWorkflowName, testKey)
	suite.Assert().Equal(testValue, savedValue)
}

func (suite *ContextConfigurationTestSuite) TestSetConfigNodeScope() {
	err := suite.ctxConfiguration.Set(testKey, testValue, NodeScope)
	suite.Require().NoError(err)
	savedValue := suite.getValueFromKVStore(suite.cfg.NATS.KeyValueStoreNodeName, testKey)
	suite.Assert().Equal(testValue, savedValue)
}

func (suite *ContextConfigurationTestSuite) TestSetConfigScopeDefault() {
	err := suite.ctxConfiguration.Set(testKey, testValue)
	suite.Require().NoError(err)
	savedValue := suite.getValueFromKVStore(suite.cfg.NATS.KeyValueStoreNodeName, testKey)
	suite.Assert().Equal(testValue, savedValue)
}

func (suite *ContextConfigurationTestSuite) TestGetConfigProjectScope() {
	// from project bucket
	err := suite.ctxConfiguration.Set(testKey, testValue, ProjectScope)
	suite.Require().NoError(err)
	savedValue, err := suite.ctxConfiguration.Get(testKey, ProjectScope)
	suite.Require().NoError(err)
	suite.Assert().Equal(testValue, savedValue)
}

func (suite *ContextConfigurationTestSuite) TestGetConfigWorkflowScope() {
	err := suite.ctxConfiguration.Set(testKey, testValue, WorkflowScope)
	suite.Require().NoError(err)
	savedValue, err := suite.ctxConfiguration.Get(testKey, WorkflowScope)
	suite.Require().NoError(err)
	suite.Assert().Equal(testValue, savedValue)
}

func (suite *ContextConfigurationTestSuite) TestGetConfigNodeScope() {
	err := suite.ctxConfiguration.Set(testKey, testValue, NodeScope)
	suite.Require().NoError(err)
	savedValue, err := suite.ctxConfiguration.Get(testKey, NodeScope)
	suite.Require().NoError(err)
	suite.Assert().Equal(testValue, savedValue)
}

func (suite *ContextConfigurationTestSuite) TestGetConfigScopeDefault() {
	nodeValue := "node"
	workflowValue := "workflow"
	projectValue := "project"

	err := suite.ctxConfiguration.Set(testKey, projectValue, ProjectScope)
	suite.Require().NoError(err)
	savedValue, err := suite.ctxConfiguration.Get(testKey)
	suite.Require().NoError(err)
	suite.Assert().Equal(projectValue, savedValue)

	err = suite.ctxConfiguration.Set(testKey, workflowValue, WorkflowScope)
	suite.Require().NoError(err)
	savedValue, err = suite.ctxConfiguration.Get(testKey)
	suite.Require().NoError(err)
	suite.Assert().Equal(workflowValue, savedValue)

	err = suite.ctxConfiguration.Set(testKey, nodeValue, NodeScope)
	suite.Require().NoError(err)
	savedValue, err = suite.ctxConfiguration.Get(testKey)
	suite.Require().NoError(err)
	suite.Assert().Equal(nodeValue, savedValue)
}

func (suite *ContextConfigurationTestSuite) TestDeleteConfigProjectScope() {
	err := suite.ctxConfiguration.Set(testKey, testValue, ProjectScope)
	suite.Require().NoError(err)
	err = suite.ctxConfiguration.Delete(testKey, ProjectScope)
	suite.Require().NoError(err)
	_, err = suite.ctxConfiguration.Get(testKey, ProjectScope)
	suite.Require().Error(err)
}

func (suite *ContextConfigurationTestSuite) TestDeleteConfigWorkflowScope() {
	err := suite.ctxConfiguration.Set(testKey, testValue, WorkflowScope)
	suite.Require().NoError(err)
	err = suite.ctxConfiguration.Delete(testKey, WorkflowScope)
	suite.Require().NoError(err)
	_, err = suite.ctxConfiguration.Get(testKey, WorkflowScope)
	suite.Require().Error(err)
}

func (suite *ContextConfigurationTestSuite) TestDeleteConfigNodeScope() {
	err := suite.ctxConfiguration.Set(testKey, testValue, NodeScope)
	suite.Require().NoError(err)
	err = suite.ctxConfiguration.Delete(testKey, NodeScope)
	suite.Require().NoError(err)
	_, err = suite.ctxConfiguration.Get(testKey, NodeScope)
	suite.Require().Error(err)
}

func (suite *ContextConfigurationTestSuite) TestDeleteConfigScopeDefault() {
	err := suite.ctxConfiguration.Set(testKey, testValue)
	suite.Require().NoError(err)
	err = suite.ctxConfiguration.Delete(testKey)
	suite.Require().NoError(err)
	_, err = suite.ctxConfiguration.Get(testKey)
	suite.Require().Error(err)
}
