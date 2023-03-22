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

type HandlerContextParams struct {
	Cfg                config.Config
	NC                 *nats.Conn
	MongoManager       mongodb.Manager
	Logger             *simplelogger.SimpleLogger
	PublishMsg         PublishMsgFunc
	PublishAny         PublishAnyFunc
	ContextObjectStore *contextObjectStore
}

type HandlerContext struct {
	cfg         config.Config
	publishMsg  PublishMsgFunc
	publishAny  PublishAnyFunc
	reqMsg      *KreNatsMessage
	Logger      *simplelogger.SimpleLogger
	ObjectStore *contextObjectStore
	Prediction  *contextPrediction
	Measurement *contextMeasurement
	DB          *contextDatabase
}

func NewHandlerContext(params *HandlerContextParams) *HandlerContext {
	return &HandlerContext{
		cfg:         params.Cfg,
		publishMsg:  params.PublishMsg,
		publishAny:  params.PublishAny,
		Logger:      params.Logger,
		ObjectStore: params.ContextObjectStore,
		Prediction:  NewContextPrediction(params.Cfg, params.NC, params.Logger),
		Measurement: NewContextMeasurement(params.Cfg, params.Logger),
		DB:          NewContextDatabase(params.Cfg, params.NC, params.MongoManager, params.Logger),
	}
}

// Path will return the relative path given as an argument as a full path.
func (c *HandlerContext) Path(relativePath string) string {
	return path.Join(c.cfg.BasePath, relativePath)
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
