package messaging

import (
	"github.com/go-logr/logr"
	kai "github.com/konstellation-io/kre-runners/go-sdk/v1/protos"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Scope string

const (
	ProductScope  Scope = "product"
	WorkflowScope Scope = "workflow"
	ProcessScope  Scope = "process"
)

type Messaging struct {
	logger         logr.Logger
	nats           *nats.Conn
	jetstream      nats.JetStreamContext
	requestMessage *kai.KaiNatsMessage
}

func NewMessaging(logger logr.Logger, nats *nats.Conn, jetstream nats.JetStreamContext,
	requestMessage *kai.KaiNatsMessage) *Messaging {
	return &Messaging{logger, nats, jetstream, requestMessage}
}

func (ms Messaging) SendOutput(response proto.Message, channelOpt ...string) error {
	return ms.publishMsg(response, ms.requestMessage.GetRequestId(), kai.MessageType_OK, ms.getOptionalString(channelOpt))
}

func (ms Messaging) SendOutputWithRequestID(response proto.Message, requestID string, channelOpt ...string) error {
	return ms.publishMsg(response, requestID, kai.MessageType_OK, ms.getOptionalString(channelOpt))
}

func (ms Messaging) SendAny(response *anypb.Any, channelOpt ...string) {
	ms.publishAny(response, ms.requestMessage.GetRequestId(), kai.MessageType_OK, ms.getOptionalString(channelOpt))
}

func (ms Messaging) SendAnyWithRequestID(response *anypb.Any, requestID string, channelOpt ...string) {
	ms.publishAny(response, requestID, kai.MessageType_OK, ms.getOptionalString(channelOpt))
}

func (ms Messaging) SendEarlyReply(response proto.Message, channelOpt ...string) error {
	return ms.publishMsg(response, ms.requestMessage.GetRequestId(), kai.MessageType_EARLY_REPLY, ms.getOptionalString(channelOpt))
}

func (ms Messaging) SendEarlyExit(response proto.Message, channelOpt ...string) error {
	return ms.publishMsg(response, ms.requestMessage.GetRequestId(), kai.MessageType_EARLY_EXIT, ms.getOptionalString(channelOpt))
}

func (ms Messaging) IsMessageOK() bool {
	return ms.requestMessage.MessageType == kai.MessageType_OK
}

func (ms Messaging) IsMessageError() bool {
	return ms.requestMessage.MessageType == kai.MessageType_ERROR
}

func (ms Messaging) IsMessageEarlyReply() bool {
	return ms.requestMessage.MessageType == kai.MessageType_EARLY_REPLY
}

func (ms Messaging) IsMessageEarlyExit() bool {
	return ms.requestMessage.MessageType == kai.MessageType_EARLY_EXIT
}
