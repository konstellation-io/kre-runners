package kre

import (
	"path"
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/konstellation-io/kre-runners/kre-go/v4/config"
	"github.com/konstellation-io/kre-runners/kre-go/v4/mongodb"
	"github.com/konstellation-io/kre/libs/simplelogger"
)

var (
	saveMetricTimeout = 1 * time.Second
	saveDataTimeout   = 1 * time.Second
	getDataTimeout    = 1 * time.Second
)

const (
	defaultValue = ""
)

type PublishMsgFunc = func(response proto.Message, reqMsg *KreNatsMessage, msgType MessageType, channel string) error
type PublishAnyFunc = func(response *anypb.Any, reqMsg *KreNatsMessage, msgType MessageType, channel string)
type StoreObjectFunc = func(key string, payload []byte) error
type GetObjectFunc = func(key string) ([]byte, error)
type DeleteObjectFunc = func(key string) error
type SaveConfigFunc = func(key, value, keyValueStore string) error
type GetConfigFunc = func(key, keyValueStore string) (string, error)

type HandlerContext struct {
	cfg          config.Config
	values       map[string]interface{}
	publishMsg   PublishMsgFunc
	publishAny   PublishAnyFunc
	storeObject  StoreObjectFunc
	getObject    GetObjectFunc
	deleteObject DeleteObjectFunc
	saveConfig   SaveConfigFunc
	getConfig    GetConfigFunc
	reqMsg       *KreNatsMessage
	Logger       *simplelogger.SimpleLogger
	Prediction   *contextPrediction
	Measurement  *contextMeasurement
	DB           *contextDatabase
}

func NewHandlerContext(
	cfg config.Config,
	nc *nats.Conn,
	mongoM mongodb.Manager,
	logger *simplelogger.SimpleLogger,
	publishMsg PublishMsgFunc,
	publishAny PublishAnyFunc,
	storeObject StoreObjectFunc,
	getObject GetObjectFunc,
	deleteObject DeleteObjectFunc,
	saveConfig SaveConfigFunc,
	getCofig GetConfigFunc,
) *HandlerContext {
	return &HandlerContext{
		cfg:          cfg,
		values:       map[string]interface{}{},
		publishMsg:   publishMsg,
		publishAny:   publishAny,
		getObject:    getObject,
		storeObject:  storeObject,
		deleteObject: deleteObject,
		saveConfig:   saveConfig,
		getConfig:    getCofig,
		Logger:       logger,
		Prediction:   NewContextPrediction(cfg, nc, logger),
		Measurement:  NewContextMeasurement(cfg, logger),
		DB:           NewContextDatabase(cfg, nc, mongoM, logger),
	}
}

// Path will return the relative path given as an argument as a full path.
func (c *HandlerContext) Path(relativePath string) string {
	return path.Join(c.cfg.BasePath, relativePath)
}

// Set will add the given key and value to the in-memory storage.
func (c *HandlerContext) Set(key string, value interface{}) {
	c.values[key] = value
}

// Get will return the value of the given key if it exists on the in-memory storage.
func (c *HandlerContext) Get(key string) interface{} {
	if val, ok := c.values[key]; ok {
		return val
	}
	c.Logger.Infof("Error getting value for key '%s' returning nil", key)
	return nil
}

// Get will return the value of the given key as a string if it exists on the in-memory storage.
func (c *HandlerContext) GetString(key string) string {
	v := c.Get(key)
	if val, ok := v.(string); ok {
		return val
	}
	c.Logger.Infof("Error getting value for key '%s' is not a string", key)
	return ""
}

// Get will return the value of the given key as an integer if it exists on the in-memory storage.
func (c *HandlerContext) GetInt(key string) int {
	v := c.Get(key)
	if val, ok := v.(int); ok {
		return val
	}
	c.Logger.Infof("Error getting value for key '%s' is not a int", key)
	return -1
}

// Get will return the value of the given key as a float if it exists on the in-memory storage.
func (c *HandlerContext) GetFloat(key string) float64 {
	v := c.Get(key)
	if val, ok := v.(float64); ok {
		return val
	}
	c.Logger.Infof("Error getting value for key '%s' is not a float64", key)
	return -1.0
}

// GetRequestID will return the payload's original request ID.
func (c *HandlerContext) GetRequestID() string {
	return c.reqMsg.RequestId
}

// SendOutput will send a desired typed proto payload to the node's subject.
// By specifying a channel, the message will be sent to that subject's subtopic.
//
// SendOutput converts the proto message into an any type. This means the following node will recieve
// an any type protobuf.
//
// GRPC requests can only be answered once. So once the entrypoint has been replied by the exitpoint,
// all following replies to the entrypoint from the same request will be ignored.
func (c *HandlerContext) SendOutput(response proto.Message, channelOpt ...string) error {
	return c.publishMsg(response, c.reqMsg, MessageType_OK, c.getOptionalString(channelOpt))
}

// SendAny will send any type of proto payload to the node's subject.
// By specifying a channel, the message will be sent to that subject's subtopic.
//
// As a difference from SendOutput, SendAny will not convert your proto structure.
//
// Use this function when you wish to simply redirect your node's payload without unpackaging.
// Once the entrypoint has been replied, all following replies to the entrypoint will be ignored.
func (c *HandlerContext) SendAny(response *anypb.Any, channelOpt ...string) {
	c.publishAny(response, c.reqMsg, MessageType_OK, c.getOptionalString(channelOpt))
}

// SendEarlyReply works as the SendOutput functionality
// with the addition of typing this message as an early reply.
func (c *HandlerContext) SendEarlyReply(response proto.Message, channelOpt ...string) error {
	return c.publishMsg(response, c.reqMsg, MessageType_EARLY_REPLY, c.getOptionalString(channelOpt))
}

// SendEarlyExit works as the SendOutput functionality
// with the addition of typing this message as an early exit.
func (c *HandlerContext) SendEarlyExit(response proto.Message, channelOpt ...string) error {
	return c.publishMsg(response, c.reqMsg, MessageType_EARLY_EXIT, c.getOptionalString(channelOpt))
}

// StoreObject stores the given payload in the Object Store with the given key as identifier
// If an objectStore name is given, a new Object Store will be created.
func (c *HandlerContext) StoreObject(key string, payload []byte) error {
	return c.storeObject(key, payload)
}

// GetObject retrieves the object stored in the given object store or the default object store for the workflow
// with the given key as identifier as a byte array.
func (c *HandlerContext) GetObject(key string) ([]byte, error) {
	return c.getObject(key)
}

// DeleteObject deletes the object stored in the given object store or the default object store for the workflow
// with the given key as identifier as a byte array.
func (c *HandlerContext) DeleteObject(key string) error {
	return c.deleteObject(key)
}

// SaveConfig Stores the given key and value to the given key-value storage,
// or the default key-value storage if not given any.
func (c *HandlerContext) SaveConfig(key, value string, keyValStoreOpt ...string) error {
	return c.saveConfig(key, value, c.getOptionalString(keyValStoreOpt))
}

// GetConfig retrieves the configuration given a key and an optional key-value storage name,
// if no key-value storage name is given it will use the default one.
func (c *HandlerContext) GetConfig(key string, keyValStoreOpt ...string) (string, error) {
	return c.getConfig(key, c.getOptionalString(keyValStoreOpt))
}

func (c *HandlerContext) getOptionalString(values []string) string {
	if len(values) > 0 {
		return values[0]
	}
	return defaultValue
}

// IsMessageOK returns true if the incoming message is of message type OK.
func (c *HandlerContext) IsMessageOK() bool {
	return c.reqMsg.MessageType == MessageType_OK
}

// IsMessageError returns true if the incoming message is of message type ERROR.
func (c *HandlerContext) IsMessageError() bool {
	return c.reqMsg.MessageType == MessageType_ERROR
}

// IsMessageEarlyReply returns true if the incoming message is of message type EARLY REPLY.
func (c *HandlerContext) IsMessageEarlyReply() bool {
	return c.reqMsg.MessageType == MessageType_EARLY_REPLY
}

// IsMessageEarlyExit returns true if the incoming message is of message type EARLY EXIT.
func (c *HandlerContext) IsMessageEarlyExit() bool {
	return c.reqMsg.MessageType == MessageType_EARLY_EXIT
}
