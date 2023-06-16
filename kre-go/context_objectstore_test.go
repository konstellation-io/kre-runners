package kre

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/testcontainers/testcontainers-go"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/konstellation-io/kre-runners/kre-go/v4/config"
	"github.com/konstellation-io/kre/libs/simplelogger"
)

type ContextObjectStoreTestSuite struct {
	suite.Suite
	logger         *simplelogger.SimpleLogger
	cfg            config.Config
	nc             *nats.Conn
	js             nats.JetStreamContext
	natsContainer  testcontainers.Container
	ctxObjectStore *contextObjectStore
}

func TestContextObjectStoreTestSuite(t *testing.T) {
	suite.Run(t, new(ContextObjectStoreTestSuite))
}

func (s *ContextObjectStoreTestSuite) SetupSuite() {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "nats:2.8.1",
		Cmd:          []string{"-js"},
		ExposedPorts: []string{"4222/tcp", "8222/tcp"},
		WaitingFor:   wait.ForLog("Server is ready"),
	}

	s.logger = simplelogger.New(simplelogger.LevelInfo)
	s.cfg = config.Config{
		NATS: config.ConfigNATS{
			ObjectStoreName: "object_store",
		},
	}

	natsContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err)

	natsEndpoint, err := natsContainer.Endpoint(ctx, "")
	s.Require().NoError(err)

	s.nc, err = nats.Connect(natsEndpoint)
	s.Require().NoError(err)

	s.js, err = s.nc.JetStream()
	s.Require().NoError(err)

	s.natsContainer = natsContainer
}

// TearDownSuite will close the mock controller, close the NATS connection and shutdown the mocked server
func (s *ContextObjectStoreTestSuite) TearDownSuite() {
	s.nc.Close()
	err := s.natsContainer.Terminate(context.Background())
	s.Require().NoError(err)
}

// SetupTest will run before each test
func (s *ContextObjectStoreTestSuite) SetupTest() {
	s.createObjectStore(s.cfg)
	objStore, err := NewContextObjectStore(s.cfg, s.logger, s.js)
	s.Require().NoError(err)
	s.ctxObjectStore = objStore
}

func (s *ContextObjectStoreTestSuite) TearDownTest() {
	s.deleteObjectStore(s.cfg)
}

func (s *ContextObjectStoreTestSuite) createObjectStore(cfg config.Config) {
	objStoreCfg := nats.ObjectStoreConfig{
		Bucket:  cfg.NATS.ObjectStoreName,
		Storage: nats.FileStorage,
	}
	s.js.CreateObjectStore(&objStoreCfg)
}

func (s *ContextObjectStoreTestSuite) deleteObjectStore(cfg config.Config) {
	s.js.DeleteObjectStore(cfg.NATS.ObjectStoreName)
}

func (s *ContextObjectStoreTestSuite) getValueFromObjectStore(bucket, key string) ([]byte, error) {
	objectStore, err := s.js.ObjectStore(bucket)
	if err != nil {
		return nil, err
	}

	entry, err := objectStore.GetBytes(key)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (s *ContextObjectStoreTestSuite) TestNewContextObjectStore() {
	// during setup test an object store is created
	s.Require().NotNil(s.ctxObjectStore.objStore)

	// if a given configuration with object store name empty, the context object store should be nil
	emptyCfg := config.Config{}
	emptyCtxObjectStore, _ := NewContextObjectStore(emptyCfg, s.logger, s.js)
	s.Require().Nil(emptyCtxObjectStore.objStore)
}

func (s *ContextObjectStoreTestSuite) TestObjectStoreSave() {
	// given
	key := "key"
	value := []byte("value")
	bucket := s.cfg.NATS.ObjectStoreName

	// when
	err := s.ctxObjectStore.Save(key, value)
	s.Require().NoError(err)

	// then
	entry, err := s.getValueFromObjectStore(bucket, key)
	s.Require().NoError(err)
	s.Require().Equal(value, entry)
}

func (s *ContextObjectStoreTestSuite) TestObjectStoreList() {
	// given
	keys := []string{"key1", "key2"}
	value := []byte("value")

	// when
	for _, key := range keys {
		err := s.ctxObjectStore.Save(key, value)
		s.Require().NoError(err)
	}

	// then
	buckets, err := s.ctxObjectStore.List()
	s.Require().NoError(err)
	s.Require().Equal(keys, buckets)
}

func (s *ContextObjectStoreTestSuite) TestObjectStoreListWithFilter() {
	// given
	keys := []string{"key1", "key2", "test1", "test2"}
	expected := []string{"test1", "test2"}
	value := []byte("value")

	// when
	for _, key := range keys {
		err := s.ctxObjectStore.Save(key, value)
		s.Require().NoError(err)
	}

	// then
	buckets, err := s.ctxObjectStore.List("test*")
	s.Require().NoError(err)
	s.Require().Equal(expected, buckets)
}

func (s *ContextObjectStoreTestSuite) TestObjectStorePurge() {
	// given
	keys := []string{"key1", "key2"}
	value := []byte("value")

	// when
	for _, key := range keys {
		err := s.ctxObjectStore.Save(key, value)
		s.Require().NoError(err)
	}

	// then
	err := s.ctxObjectStore.Purge()
	s.Require().NoError(err)

	buckets, err := s.ctxObjectStore.List()
	s.Require().Error(err)
	s.Require().Empty(buckets)
}

func (s *ContextObjectStoreTestSuite) TestObjectStorePurgeWithFilter() {
	// given
	keys := []string{"key1", "key2", "test1", "test2"}
	expected := []string{"key1", "key2"}
	value := []byte("value")

	// when
	for _, key := range keys {
		err := s.ctxObjectStore.Save(key, value)
		s.Require().NoError(err)
	}

	// then
	err := s.ctxObjectStore.Purge("test*")
	s.Require().NoError(err)

	buckets, err := s.ctxObjectStore.List()
	s.Require().NoError(err)
	s.Require().Equal(expected, buckets)
}

func (s *ContextObjectStoreTestSuite) TestObjectStoreSaveWithEmptyKey() {
	// given
	key := ""
	value := []byte("value")

	// when
	err := s.ctxObjectStore.Save(key, value)

	// then
	s.Require().Error(err)
}

func (s *ContextObjectStoreTestSuite) TestObjectStoreSaveWithEmptyValue() {
	// given
	key := "key"
	value := []byte("")

	// when
	err := s.ctxObjectStore.Save(key, value)

	// then
	s.Require().NoError(err)

	// additionaly an error returns when trying to get the value
	entry, err := s.getValueFromObjectStore(s.cfg.NATS.ObjectStoreName, key)
	s.Require().NoError(err)
	s.Require().Empty(entry)
}

func (s *ContextObjectStoreTestSuite) TestObjectStoreSaveWithNilValue() {
	// given
	key := "key"

	// when
	err := s.ctxObjectStore.Save(key, nil)

	// then
	s.Require().Error(err)

	// additionaly an error returns when trying to get the value
	entry, err := s.getValueFromObjectStore(s.cfg.NATS.ObjectStoreName, key)
	s.Require().Error(err)
	s.Require().Empty(entry)
}

func (s *ContextObjectStoreTestSuite) TestObjectStoreGet() {
	// given
	key := "key"
	value := []byte("value")

	// when
	err := s.ctxObjectStore.Save(key, value)
	s.Require().NoError(err)

	entry, err := s.ctxObjectStore.Get(key)
	s.Require().NoError(err)

	// then
	s.Require().Equal(value, entry)
}

func (s *ContextObjectStoreTestSuite) TestObjectStoreGetWithEmptyKey() {
	// given
	key := ""

	// when
	entry, err := s.ctxObjectStore.Get(key)

	// then
	s.Require().Error(err)
	s.Require().Nil(entry)
}

func (s *ContextObjectStoreTestSuite) TestObjectStoreGetWithNonExistingKey() {
	// given
	key := "non_existing_key"

	// when
	entry, err := s.ctxObjectStore.Get(key)

	// then
	s.Require().Error(err)
	s.Require().Nil(entry)
}

func (s *ContextObjectStoreTestSuite) TestObjectStoreDelete() {
	// given
	key := "key"
	value := []byte("value")

	// given
	entry, err := s.ctxObjectStore.Get(key)
	s.Require().Error(err)
	s.Assert().Empty(entry)

	// when
	err = s.ctxObjectStore.Save(key, value)
	s.Require().NoError(err)

	err = s.ctxObjectStore.Delete(key)
	s.Require().NoError(err)

	entry, err = s.ctxObjectStore.Get(key)

	// then
	s.Assert().Error(err)
	s.Assert().Empty(entry)
}

func (s *ContextObjectStoreTestSuite) TestObjectStoreDeleteWithEmptyKey() {
	// given
	key := ""

	// when
	err := s.ctxObjectStore.Delete(key)

	// then
	s.Require().Error(err)
}

func (s *ContextObjectStoreTestSuite) TestObjectStoreDeleteWithNonExistingKey() {
	// given
	key := "non_existing_key"

	// when
	err := s.ctxObjectStore.Delete(key)

	// then
	s.Require().Error(err)
}
