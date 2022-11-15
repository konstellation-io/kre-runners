package kre

import (
	"errors"
	"fmt"
	"time"

	"github.com/konstellation-io/kre/libs/simplelogger"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/konstellation-io/kre-runners/kre-go/config"
	"github.com/konstellation-io/kre-runners/kre-go/mongodb"
)

const (
	MessageThreshold = 1024 * 1024
)

var ErrMessageToBig = errors.New("compressed message exceeds maximum size allowed of 1 MB")
var ErrMsgAck = "Error in message ack: %s"

type Runner struct {
	logger         *simplelogger.SimpleLogger
	cfg            config.Config
	nc             *nats.Conn
	js             nats.JetStreamContext
	handlerContext *HandlerContext
	handlerManager *HandlerManager
}

// NewRunner creates a new Runner instance, initializing a new handler context within and runs
// the given handler init func.
func NewRunner(
	logger *simplelogger.SimpleLogger,
	cfg config.Config,
	nc *nats.Conn,
	js nats.JetStreamContext,
	handlerManager *HandlerManager,
	handlerInit HandlerInit,
	mongoM *mongodb.MongoDB,
) *Runner {
	runner := &Runner{
		logger:         logger,
		cfg:            cfg,
		nc:             nc,
		js:             js,
		handlerManager: handlerManager,
	}

	c := NewHandlerContext(cfg, nc, mongoM, logger, runner.publishMsg, runner.publishAny)
	handlerInit(c)

	runner.handlerContext = c

	return runner
}

// ProcessMessage parses the incoming NATS message and executes the appropiate handler function
// taking into account the origin's node of the message.
func (r *Runner) ProcessMessage(msg *nats.Msg) {
	var (
		start = time.Now().UTC()
	)

	requestMsg, err := r.newRequestMessage(msg.Data)
	if err != nil {
		errMsg := fmt.Sprintf("Error parsing msg.data coming from subject %s because is not a valid protobuf: %s", msg.Subject, err)
		r.processRunnerError(msg, errMsg, requestMsg.RequestId, start, requestMsg.FromNode)
		return
	}

	r.logger.Infof("Received a message from '%s' with requestId '%s'", msg.Subject, requestMsg.RequestId)

	// Make a shallow copy of the ctx object to set inside the request msg.
	hCtx := r.handlerContext
	hCtx.reqMsg = requestMsg

	handler := r.handlerManager.GetHandler(requestMsg.FromNode)
	if handler == nil {
		errMsg := fmt.Sprintf("Error missing handler for node '%s'", requestMsg.FromNode)
		r.processRunnerError(msg, errMsg, requestMsg.RequestId, start, requestMsg.FromNode)
		return
	}

	err = handler(hCtx, requestMsg.Payload)
	if err != nil {
		errMsg := fmt.Sprintf("Error in node '%s' executing handler for node '%s': %s", r.cfg.NodeName, requestMsg.FromNode, err)
		r.processRunnerError(msg, errMsg, requestMsg.RequestId, start, requestMsg.FromNode)
		return
	}

	// Tell NATS we don't need to receive the message anymore and we are done processing it.
	ackErr := msg.Ack()
	if ackErr != nil {
		r.logger.Errorf(ErrMsgAck, ackErr)
	}

	end := time.Now().UTC()
	r.saveElapsedTime(start, end, requestMsg.FromNode, true)
}

func (r *Runner) processRunnerError(msg *nats.Msg, errMsg string, requestID string, start time.Time, fromNode string) {
	ackErr := msg.Ack()
	if ackErr != nil {
		r.logger.Errorf(ErrMsgAck, ackErr)
	}

	r.logger.Error(errMsg)
	r.publishError(requestID, errMsg)

	end := time.Now().UTC()
	r.saveElapsedTime(start, end, fromNode, false)
}

func (r *Runner) newRequestMessage(data []byte) (*KreNatsMessage, error) {
	requestMsg := &KreNatsMessage{}

	var err error
	if r.isCompressed(data) {
		data, err = r.uncompressData(data)
		if err != nil {
			r.logger.Errorf("error reading compressed message: %s", err)
			return nil, err
		}
	}

	err = proto.Unmarshal(data, requestMsg)

	return requestMsg, err
}

// publishMsg will send a desired payload to the node's output subject.
func (r *Runner) publishMsg(msg proto.Message, reqMsg *KreNatsMessage, msgType MessageType, channel string) error {
	payload, err := anypb.New(msg)
	if err != nil {
		return fmt.Errorf("the handler result is not a valid protobuf: %s", err)
	}
	responseMsg := r.newResponseMsg(payload, reqMsg, msgType)

	r.publishResponse(responseMsg, channel)

	return nil
}

func (r *Runner) publishAny(payload *anypb.Any, reqMsg *KreNatsMessage, msgType MessageType, channel string) {
	responseMsg := r.newResponseMsg(payload, reqMsg, msgType)
	r.publishResponse(responseMsg, channel)
}

func (r *Runner) publishError(requestID, errMsg string) {
	responseMsg := &KreNatsMessage{
		RequestId:   requestID,
		Error:       errMsg,
		FromNode:    r.cfg.NodeName,
		MessageType: MessageType_ERROR,
	}
	r.publishResponse(responseMsg, "")
}

// newResponseMsg creates a KreNatsMessage that keeps previous request ID plus adding the payload we wish to send.
func (r *Runner) newResponseMsg(payload *anypb.Any, requestMsg *KreNatsMessage, msgType MessageType) *KreNatsMessage {
	return &KreNatsMessage{
		RequestId:   requestMsg.RequestId,
		Payload:     payload,
		FromNode:    r.cfg.NodeName,
		MessageType: msgType,
	}
}

func (r *Runner) publishResponse(responseMsg *KreNatsMessage, channel string) {
	outputSubject := r.getOutputSubject(channel)

	outputMsg, err := proto.Marshal(responseMsg)
	if err != nil {
		r.logger.Errorf("Error generating output result because handler result is not a serializable Protobuf: %s", err)
		return
	}

	outputMsg, err = r.prepareOutputMessage(outputMsg)
	if err != nil {
		r.logger.Errorf("Error preparing output msg: %s", err)
		return
	}

	r.logger.Infof("Publishing response to '%s' subject", outputSubject)

	_, err = r.js.Publish(outputSubject, outputMsg)
	if err != nil {
		r.logger.Errorf("Error publishing output: %s", err)
	}
}

func (r *Runner) getOutputSubject(channel string) string {
	outputSubject := r.cfg.NATS.OutputSubject
	if channel != "" {
		return fmt.Sprintf("%s.%s", outputSubject, channel)
	}
	return outputSubject
}

// prepareOutputMessage will check the length of the message and compress it if necessary.
// Fails on compressed messages bigger than the threshold.
func (r *Runner) prepareOutputMessage(msg []byte) ([]byte, error) {
	if len(msg) <= MessageThreshold {
		return msg, nil
	}

	outMsg, err := r.compressData(msg)
	if err != nil {
		return nil, err
	}

	if len(outMsg) > MessageThreshold {
		return nil, ErrMessageToBig
	}

	r.logger.Infof("Original message size: %s. Compressed: %s", sizeInKB(msg), sizeInKB(outMsg))

	return outMsg, nil
}

// saveElapsedTime stores in InfluxDB how much time did it take the node to run the handler,
// also saves if the request was succesfully processed.
func (r *Runner) saveElapsedTime(start time.Time, end time.Time, fromNode string, success bool) {
	elapsed := end.Sub(start)

	tags := map[string]string{
		"from_node": fromNode,
	}

	fields := map[string]interface{}{
		"elapsed_ms": elapsed.Seconds() * 1000,
		"success":    success,
	}

	r.handlerContext.Measurement.Save("node_elapsed_time", fields, tags)
}

func sizeInKB(s []byte) string {
	return fmt.Sprintf("%.2f KB", float32(len(s))/1024)
}
