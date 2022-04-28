package kre

import (
	"errors"
	"path"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/konstellation-io/kre/libs/simplelogger"
	"github.com/nats-io/nats.go"

	"github.com/konstellation-io/kre-runners/kre-go/config"
	"github.com/konstellation-io/kre-runners/kre-go/mongodb"
)

var (
	saveMetricTimeout = 1 * time.Second
	saveDataTimeout   = 1 * time.Second
	getDataTimeout    = 1 * time.Second
)

type ReplyFunc = func(subject string, response proto.Message) error
type SendOutputFunc = func(subject string, entrypointSubject string, response proto.Message) error

type HandlerContext struct {
	cfg         config.Config
	values      map[string]interface{}
	reply       ReplyFunc
	sendOutput  SendOutputFunc
	reqMsg      *KreNatsMessage
	Logger      *simplelogger.SimpleLogger
	Prediction  *contextPrediction
	Measurement *contextMeasurement
	DB          *contextData
}

func NewHandlerContext(cfg config.Config, nc *nats.Conn, mongoM mongodb.Manager, logger *simplelogger.SimpleLogger, reply ReplyFunc, sendOutput SendOutputFunc) *HandlerContext {
	return &HandlerContext{
		cfg:        cfg,
		values:     map[string]interface{}{},
		reply:      reply,
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

// Reply sends a reply to the entrypoint. The workflow execution continues.
// Use this function when you need to reply faster than the workflow execution duration.
func (c *HandlerContext) Reply(response proto.Message) error {
	if c.reqMsg.Replied {
		return errors.New("error the message was replied previously")
	}

	c.reqMsg.Replied = true
	return c.reply(c.reqMsg.Reply, response)
}

// Reply sends a reply to the entrypoint. The workflow execution continues.
// Use this function when you need to reply faster than the workflow execution duration.
func (c *HandlerContext) SendOutput(response proto.Message) error {

	return c.sendOutput(c.cfg.NATS.OutputSubject, c.reqMsg.Reply, response)
}
