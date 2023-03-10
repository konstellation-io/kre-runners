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
type SetConfigFunc = func(key, value string, scope Scope) error
type GetConfigFunc = func(key string, scopes []Scope) (string, error)
type DeleteConfigFunc = func(key string, scope Scope) error

type Scope string

const (
	ScopeProject  Scope = "project"
	ScopeWorkflow Scope = "workflow"
	ScopeNode     Scope = "node"
)

type HandlerContext struct {
	cfg          config.Config
	values       map[string]interface{}
	publishMsg   PublishMsgFunc
	publishAny   PublishAnyFunc
	storeObject  StoreObjectFunc
	getObject    GetObjectFunc
	deleteObject DeleteObjectFunc
	setConfig    SetConfigFunc
	getConfig    GetConfigFunc
	deleteConfig DeleteConfigFunc
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
	setConfig SetConfigFunc,
	getConfig GetConfigFunc,
	deleteConfig DeleteConfigFunc,
) *HandlerContext {
	return &HandlerContext{
		cfg:          cfg,
		values:       map[string]interface{}{},
		publishMsg:   publishMsg,
		publishAny:   publishAny,
		getObject:    getObject,
		storeObject:  storeObject,
		deleteObject: deleteObject,
		setConfig:    setConfig,
		getConfig:    getConfig,
		deleteConfig: deleteConfig,
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

// GetString will return the value of the given key as a string if it exists on the in-memory storage.
func (c *HandlerContext) GetString(key string) string {
	v := c.Get(key)
	if val, ok := v.(string); ok {
		return val
	}
	c.Logger.Infof("Error getting value for key '%s' is not a string", key)
	return ""
}

// GetInt will return the value of the given key as an integer if it exists on the in-memory storage.
func (c *HandlerContext) GetInt(key string) int {
	v := c.Get(key)
	if val, ok := v.(int); ok {
		return val
	}
	c.Logger.Infof("Error getting value for key '%s' is not a int", key)
	return -1
}

// GetFloat will return the value of the given key as a float if it exists on the in-memory storage.
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

// SetConfig set the given key and value to an optional scoped key-value storage,
// or the default key-value storage (Node's) if not given any.
func (c *HandlerContext) SetConfig(key, value string, scopeOpt ...Scope) error {
	return c.setConfig(key, value, c.getOptionalScope(scopeOpt))
}

// GetConfig retrieves the configuration given a key from an optional scoped key-value storage,
// if no scoped key-value storage is given it will search in all the scopes starting by Node then upwards.
func (c *HandlerContext) GetConfig(key string, scopeOpt ...Scope) (string, error) {
	return c.getConfig(key, scopeOpt)
}

// DeleteConfig retrieves the configuration given a key from an optional scoped key-value storage,
// if no key-value storage is given it will use the default one (Node's).
func (c *HandlerContext) DeleteConfig(key string, scopeOpt ...Scope) error {
	return c.deleteConfig(key, c.getOptionalScope(scopeOpt))
}

func (c *HandlerContext) getOptionalString(values []string) string {
	if len(values) > 0 {
		return values[0]
	}
	return defaultValue
}

func (c *HandlerContext) getOptionalScope(scopes []Scope) Scope {
	if len(scopes) > 0 {
		return scopes[0]
	}
	return ScopeNode
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
