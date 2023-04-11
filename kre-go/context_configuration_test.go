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
	testKey   = "key"
	testValue = "value"
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

	suite.createKVStores(suite.cfg)

	suite.ctxConfiguration, err = NewContextConfiguration(suite.cfg, suite.logger, suite.js)
	suite.Require().NoError(err)
}

func (suite *ContextConfigurationTestSuite) TearDownTest() {
	err := suite.js.DeleteKeyValue(suite.cfg.NATS.KeyValueStoreProjectName)
	suite.Require().NoError(err)

	err = suite.js.DeleteKeyValue(suite.cfg.NATS.KeyValueStoreWorkflowName)
	suite.Require().NoError(err)

	err = suite.js.DeleteKeyValue(suite.cfg.NATS.KeyValueStoreNodeName)
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

func (suite *ContextConfigurationTestSuite) TestSetConfigScopeProject() {
	err := suite.ctxConfiguration.Set(testKey, testValue, ScopeProject)
	suite.Require().NoError(err)
	savedValue := suite.getValueFromKVStore(suite.cfg.NATS.KeyValueStoreProjectName, testKey)
	suite.Assert().Equal(testValue, savedValue)
}

func (suite *ContextConfigurationTestSuite) TestSetConfigScopeWorkflow() {
	err := suite.ctxConfiguration.Set(testKey, testValue, ScopeWorkflow)
	suite.Require().NoError(err)
	savedValue := suite.getValueFromKVStore(suite.cfg.NATS.KeyValueStoreWorkflowName, testKey)
	suite.Assert().Equal(testValue, savedValue)
}

func (suite *ContextConfigurationTestSuite) TestSetConfigScopeNode() {
	err := suite.ctxConfiguration.Set(testKey, testValue, ScopeNode)
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

func (suite *ContextConfigurationTestSuite) TestGetConfigScopeProject() {
	// from project bucket
	err := suite.ctxConfiguration.Set(testKey, testValue, ScopeProject)
	suite.Require().NoError(err)
	savedValue, err := suite.ctxConfiguration.Get(testKey, ScopeProject)
	suite.Require().NoError(err)
	suite.Assert().Equal(testValue, savedValue)
}

func (suite *ContextConfigurationTestSuite) TestGetConfigScopeWorkflow() {
	err := suite.ctxConfiguration.Set(testKey, testValue, ScopeWorkflow)
	suite.Require().NoError(err)
	savedValue, err := suite.ctxConfiguration.Get(testKey, ScopeWorkflow)
	suite.Require().NoError(err)
	suite.Assert().Equal(testValue, savedValue)
}

func (suite *ContextConfigurationTestSuite) TestGetConfigScopeNode() {
	err := suite.ctxConfiguration.Set(testKey, testValue, ScopeNode)
	suite.Require().NoError(err)
	savedValue, err := suite.ctxConfiguration.Get(testKey, ScopeNode)
	suite.Require().NoError(err)
	suite.Assert().Equal(testValue, savedValue)
}

func (suite *ContextConfigurationTestSuite) TestGetConfigScopeDefault() {
	nodeValue := "node"
	workflowValue := "workflow"
	projectValue := "project"

	err := suite.ctxConfiguration.Set(testKey, projectValue, ScopeProject)
	suite.Require().NoError(err)
	savedValue, err := suite.ctxConfiguration.Get(testKey)
	suite.Require().NoError(err)
	suite.Assert().Equal(projectValue, savedValue)

	err = suite.ctxConfiguration.Set(testKey, workflowValue, ScopeWorkflow)
	suite.Require().NoError(err)
	savedValue, err = suite.ctxConfiguration.Get(testKey)
	suite.Require().NoError(err)
	suite.Assert().Equal(workflowValue, savedValue)

	err = suite.ctxConfiguration.Set(testKey, nodeValue, ScopeNode)
	suite.Require().NoError(err)
	savedValue, err = suite.ctxConfiguration.Get(testKey)
	suite.Require().NoError(err)
	suite.Assert().Equal(nodeValue, savedValue)
}

func (suite *ContextConfigurationTestSuite) TestDeleteConfigScopeProject() {
	err := suite.ctxConfiguration.Set(testKey, testValue, ScopeProject)
	suite.Require().NoError(err)
	err = suite.ctxConfiguration.Delete(testKey, ScopeProject)
	suite.Require().NoError(err)
	_, err = suite.ctxConfiguration.Get(testKey, ScopeProject)
	suite.Require().Error(err)
}

func (suite *ContextConfigurationTestSuite) TestDeleteConfigScopeWorkflow() {
	err := suite.ctxConfiguration.Set(testKey, testValue, ScopeWorkflow)
	suite.Require().NoError(err)
	err = suite.ctxConfiguration.Delete(testKey, ScopeWorkflow)
	suite.Require().NoError(err)
	_, err = suite.ctxConfiguration.Get(testKey, ScopeWorkflow)
	suite.Require().Error(err)
}

func (suite *ContextConfigurationTestSuite) TestDeleteConfigScopeNode() {
	err := suite.ctxConfiguration.Set(testKey, testValue, ScopeNode)
	suite.Require().NoError(err)
	err = suite.ctxConfiguration.Delete(testKey, ScopeNode)
	suite.Require().NoError(err)
	_, err = suite.ctxConfiguration.Get(testKey, ScopeNode)
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
