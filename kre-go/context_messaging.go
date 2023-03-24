package kre

import (
	"fmt"
	"github.com/konstellation-io/kre-runners/kre-go/v4/config"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type ContextMessaging struct {
	cfg             config.Config
	messagingClient MessagingClient
	requestId       string
}

func NewContextMessaging(cfg config.Config, client MessagingClient) *ContextMessaging {
	return &ContextMessaging{
		cfg:             cfg,
		messagingClient: client,
	}
}

func (m *ContextMessaging) SetRequestId(requestId string) {
	m.requestId = requestId
}

func (m *ContextMessaging) SendOutput(response proto.Message, channelOpt ...string) error {
	return m.messagingClient.PublishResponse(response, MessageType_OK, m.getOptionalString(channelOpt))
}

// publishMsg will send a desired payload to the node's output subject.
func (m *ContextMessaging) publishMsg(msg proto.Message, msgType MessageType, channel string) error {
	payload, err := anypb.New(msg)
	if err != nil {
		return fmt.Errorf("the handler result is not a valid protobuf: %s", err)
	}
	responseMsg := m.newResponseMsg(payload, msgType)

	r.messagingClient.PublishResponse(responseMsg, channel)

	return nil
}

// newResponseMsg creates a KreNatsMessage that keeps previous request ID plus adding the payload we wish to send.
func (m *ContextMessaging) newResponseMsg(payload *anypb.Any, msgType MessageType) *KreNatsMessage {
	return &KreNatsMessage{
		RequestId:   m.requestId,
		Payload:     payload,
		FromNode:    m.cfg.NodeName,
		MessageType: msgType,
	}
}

func (m *ContextMessaging) getOptionalString(values []string) string {
	if len(values) > 0 {
		return values[0]
	}
	return defaultValue
}
