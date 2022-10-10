package kre

import (
	"path"
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/konstellation-io/kre-runners/kre-go/config"
	"github.com/konstellation-io/kre-runners/kre-go/mongodb"
	"github.com/konstellation-io/kre/libs/simplelogger"
)

var (
	saveMetricTimeout = 1 * time.Second
	saveDataTimeout   = 1 * time.Second
	getDataTimeout    = 1 * time.Second
)

const (
	defaultChannel = ""
)

type PublishMsgFunc = func(response proto.Message, reqMsg *KreNatsMessage, msgType MessageType, channel string) error
type PublishAnyFunc = func(response *anypb.Any, reqMsg *KreNatsMessage, msgType MessageType, channel string) error

type HandlerContext struct {
	cfg         config.Config
	values      map[string]interface{}
	publishMsg  PublishMsgFunc
	publishAny  PublishAnyFunc
	reqMsg      *KreNatsMessage
	Logger      *simplelogger.SimpleLogger
	Prediction  *contextPrediction
	Measurement *contextMeasurement
	DB          *contextData
}

func NewHandlerContext(cfg config.Config, nc *nats.Conn, mongoM mongodb.Manager,
	logger *simplelogger.SimpleLogger, publishMsg PublishMsgFunc, publishAny PublishAnyFunc) *HandlerContext {
	return &HandlerContext{
		cfg:        cfg,
		values:     map[string]interface{}{},
		publishMsg: publishMsg,
		publishAny: publishAny,
		Logger:     logger,
		Prediction: &contextPrediction{
			cfg:    cfg,
			nc:     nc,
			logger: logger,
		},
		Measurement: NewContextMeasurement(cfg, logger),
		DB: &contextData{
			cfg:    cfg,
			nc:     nc,
			mongoM: mongoM,
			logger: logger,
		},
	}
}

func (c *HandlerContext) Path(relativePath string) string {
	return path.Join(c.cfg.BasePath, relativePath)
}

func (c *HandlerContext) Set(key string, value interface{}) {
	c.values[key] = value
}

func (c *HandlerContext) Get(key string) interface{} {
	if val, ok := c.values[key]; ok {
		return val
	}
	c.Logger.Infof("Error getting value for key '%s' returning nil", key)
	return nil
}

func (c *HandlerContext) GetString(key string) string {
	v := c.Get(key)
	if val, ok := v.(string); ok {
		return val
	}
	c.Logger.Infof("Error getting value for key '%s' is not a string", key)
	return ""
}

func (c *HandlerContext) GetInt(key string) int {
	v := c.Get(key)
	if val, ok := v.(int); ok {
		return val
	}
	c.Logger.Infof("Error getting value for key '%s' is not a int", key)
	return -1
}

func (c *HandlerContext) GetFloat(key string) float64 {
	v := c.Get(key)
	if val, ok := v.(float64); ok {
		return val
	}
	c.Logger.Infof("Error getting value for key '%s' is not a float64", key)
	return -1.0
}

// GetRequestID will return the message's request ID
func (c *HandlerContext) GetRequestID() string {
	return c.reqMsg.RequestId
}

// SendOutput will send a desired payload to the node's subject.
// Once the entrypoint has been replied, all following replies to the entrypoint will be ignored.
func (c *HandlerContext) SendOutput(response proto.Message, channelOpt ...string) error {
	channel := defaultChannel
	if len(channelOpt) > 0 {
		channel = channelOpt[0]
	}
	return c.publishMsg(response, c.reqMsg, MessageType_OK, channel)
}

// SendEarlyReply publishes the desired response to this node's subject.
// With the addition of typing this message as an early reply.
// Use this function when you need to reply faster than the workflow execution duration.
func (c *HandlerContext) SendEarlyReply(response proto.Message) error {
	return c.publishMsg(response, c.reqMsg, MessageType_EARLY_REPLY, defaultChannel)
}

// SendEarlyExit publishes the desired response to this node's subject.
// With the addition of typing this message as an early exit.
// Use this function when you want to report a custom error in your workflow execution.
func (c *HandlerContext) SendEarlyExit(response proto.Message) error {
	return c.publishMsg(response, c.reqMsg, MessageType_EARLY_EXIT, defaultChannel)
}

// GetRequestMessageType returns the message type of the incoming request
func (c *HandlerContext) GetRequestMessageType() MessageType {
	return c.reqMsg.MessageType
}

func (c *HandlerContext) SendAny(response *anypb.Any, channelOpt ...string) error {
	channel := defaultChannel
	if len(channelOpt) > 0 {
		channel = channelOpt[0]
	}
	return c.publishAny(response, c.reqMsg, MessageType_OK, channel)
}
