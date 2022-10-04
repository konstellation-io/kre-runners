package kre

import (
	"errors"
	"path"
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	"github.com/konstellation-io/kre-runners/kre-go/config"
	"github.com/konstellation-io/kre-runners/kre-go/mongodb"
	"github.com/konstellation-io/kre/libs/simplelogger"
)

var (
	saveMetricTimeout = 1 * time.Second
	saveDataTimeout   = 1 * time.Second
	getDataTimeout    = 1 * time.Second
)

type EarlyReplyFunc = func(response proto.Message, requestID string) error
type SendOutputFunc = func(response proto.Message, reqMsg *KreNatsMessage) error

type HandlerContext struct {
	cfg         config.Config
	values      map[string]interface{}
	earlyReply  EarlyReplyFunc
	sendOutput  SendOutputFunc
	ReqMsg      *KreNatsMessage
	Logger      *simplelogger.SimpleLogger
	Prediction  *contextPrediction
	Measurement *contextMeasurement
	DB          *contextData
}

func NewHandlerContext(cfg config.Config, nc *nats.Conn, mongoM mongodb.Manager,
	logger *simplelogger.SimpleLogger, earlyReply EarlyReplyFunc, sendOutput SendOutputFunc) *HandlerContext {
	return &HandlerContext{
		cfg:        cfg,
		values:     map[string]interface{}{},
		earlyReply: earlyReply,
		sendOutput: sendOutput,
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

func (c *HandlerContext) GetRequestID() string {
	return c.ReqMsg.RequestId
}

// EarlyReply sends a reply to the entrypoint. The workflow execution continues.
// Use this function when you need to reply faster than the workflow execution duration.
func (c *HandlerContext) EarlyReply(response proto.Message) error {
	if c.ReqMsg.Replied {
		return errors.New("error the message was replied previously")
	}

	c.ReqMsg.Replied = true
	return c.earlyReply(response, c.ReqMsg.RequestId)
}

// SetEarlyExit changes this node's next response recipient to the entrypoint.
// Use this function when you want to halt your workflow execution.
func (c *HandlerContext) SetEarlyExit() {
	c.ReqMsg.EarlyExit = true
}

// SendOutput will send a desired payload to the node's subject.
// Once the entrypoint has been replied, all following replies to the entrypoint will be ignored.
func (c *HandlerContext) SendOutput(response proto.Message) error {
	return c.sendOutput(response, c.ReqMsg)
}

// SendEarlyExit combines both SetEarlyExit and SendOutput functions.
// As a result when invoquing this function, a response will be sent directly to the entrypoint and
// the workflow execution will be halted.
func (c *HandlerContext) SendEarlyExit(response proto.Message) error {
	c.SetEarlyExit()
	return c.sendOutput(response, c.ReqMsg)
}
