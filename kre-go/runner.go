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
	ISO8601          = "2006-01-02T15:04:05.000000"
	MessageThreshold = 1024 * 1024
)

var ErrMessageToBig = errors.New("compressed message exceeds maximum size allowed of 1 MB")

type Runner struct {
	logger         *simplelogger.SimpleLogger
	cfg            config.Config
	nc             *nats.Conn
	js             nats.JetStreamContext
	handlerContext *HandlerContext
	handlerManager *Manager
}

// NewRunner creates a new Runner instance.
func NewRunner(logger *simplelogger.SimpleLogger, cfg config.Config, nc *nats.Conn, js nats.JetStreamContext,
	handlerManager *Manager, handlerInit HandlerInit, mongoM *mongodb.MongoDB) *Runner {
	runner := &Runner{
		logger:         logger,
		cfg:            cfg,
		nc:             nc,
		js:             js,
		handlerManager: handlerManager,
	}

	// Create handler context
	c := NewHandlerContext(cfg, nc, mongoM, logger, runner.publishMsg)
	handlerInit(c)

	runner.handlerContext = c

	return runner
}

// ProcessMessage parses the incoming NATS message, executes the handler function and publishes
// the handler result to the output subject.
func (r *Runner) ProcessMessage(msg *nats.Msg) {
	start := time.Now().UTC()

	// Parse incoming message
	requestMsg, err := r.newRequestMessage(msg.Data)
	if err != nil {
		r.logger.Errorf("Error parsing msg.data because is not a valid protobuf: %s", err)
		return
	}

	r.logger.Infof("Received a message on '%s' to be published in '%s' with requestId '%s'", msg.Subject, r.cfg.NATS.OutputSubject, requestMsg.RequestId)

	if requestMsg.MessageType == MessageType_EARLY_REPLY {
		if !r.cfg.IsExitpoint { // ignore if I am not exitpoint
			return
		}
		// reroute to entrypoint if I am exitpoint
		err := r.publishMsg(requestMsg.Payload, requestMsg, MessageType_OK)
		if err != nil {
			r.logger.Errorf("error publishing early reply message: %s", err.Error())
		}

		ackErr := msg.Ack()
		if ackErr != nil {
			r.logger.Errorf("Error in message ack: %s", ackErr.Error())
		}
		return
	}

	// Make a shallow copy of the ctx object to set inside the request msg and set it to this runner.
	hCtx := r.handlerContext
	hCtx.reqMsg = requestMsg

	// Execute the handler function sending context and the payload.
	handler := r.handlerManager.GetHandler(requestMsg.FromNode)
	if handler == nil {
		errMsg := fmt.Sprintf("Error missing handler for node '%s'", requestMsg.FromNode)
		r.logger.Errorf(errMsg)
		r.publishError(requestMsg.RequestId, errMsg)

		ackErr := msg.Ack()
		if ackErr != nil {
			r.logger.Errorf("Error in message ack: %s", ackErr.Error())
		}
		return
	}
	err = handler(hCtx, requestMsg.Payload)
	// Tell NATS we don't need to receive the message anymore and we are done processing it.
	ackErr := msg.Ack()
	if ackErr != nil {
		r.logger.Errorf("Error in message ack: %s", ackErr.Error())
	}
	if err != nil {
		errMsg := fmt.Sprintf("Error in node '%s' executing handler for node '%s': %s", r.cfg.NodeName, requestMsg.FromNode, err)
		r.logger.Errorf(errMsg)
		r.publishError(requestMsg.RequestId, errMsg)
		return
	}

	end := time.Now().UTC()

	// Save the elapsed time for this node and for the workflow if it is the last node.
	r.saveElapsedTime(requestMsg, start, end, r.cfg.IsExitpoint)
}

// publishMsg will send a desired payload to the node's output subject.
func (r *Runner) publishMsg(msg proto.Message, reqMsg *KreNatsMessage, msgType MessageType) error {
	// Generate a KreNatsMessage response.
	end := time.Now().UTC()
	responseMsg, err := r.newResponseMsg(msg, reqMsg, end, msgType)
	if err != nil {
		return err
	}

	// Publish the response message to the output subject.
	r.publishResponse(responseMsg)

	return nil
}

// publishError will send a custom error to the node's output subject.
func (r *Runner) publishError(requestID, errMsg string) {
	// Generate a KreNatsMessage response.
	responseMsg := &KreNatsMessage{
		RequestId:   requestID,
		Error:       errMsg,
		MessageType: MessageType_ERROR,
		FromNode:    r.cfg.NodeName,
	}

	// Publish the response message to the output subject.
	r.publishResponse(responseMsg)
}

// newRequestMessage creates an instance of KreNatsMessage for the input string. decompress if necessary
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

// newResponseMsg creates a KreNatsMessage maintaining the tracking ID and adding the
// handler result and the tracking information for this node.
func (r *Runner) newResponseMsg(msg proto.Message, requestMsg *KreNatsMessage, end time.Time, msgType MessageType) (*KreNatsMessage, error) {
	payload, err := anypb.New(msg)
	if err != nil {
		return nil, fmt.Errorf("the handler result is not a valid protobuf: %w", err)
	}

	tracking := append(requestMsg.Tracking, &KreNatsMessage_Tracking{
		// Start time is only needed from the entrypoint node.
		NodeName: r.cfg.NodeName,
		End:      end.Format(ISO8601),
	})

	responseMsg := &KreNatsMessage{
		Replied:     requestMsg.Replied,
		TrackingId:  requestMsg.TrackingId,
		Tracking:    tracking,
		Payload:     payload,
		RequestId:   requestMsg.RequestId,
		FromNode:    r.cfg.NodeName,
		MessageType: msgType,
	}

	return responseMsg, nil
}

// prepareOutputMessage check the length of the message and compress if necessary.
// fails on compressed messages bigger than the threshold.
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

// publishResponse publishes the response in the NATS output subject.
func (r *Runner) publishResponse(responseMsg *KreNatsMessage) {
	outputSubject := r.cfg.NATS.OutputSubject

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

// saveElapsedTime stores in InfluxDB the elapsed time for the current node and the total elapsed time of the
// complete workflow if it is the last node.
func (r *Runner) saveElapsedTime(reqMsg *KreNatsMessage, start time.Time, end time.Time, isExitpoint bool) {
	prev := reqMsg.Tracking[len(reqMsg.Tracking)-1]
	prevEnd, err := time.Parse(ISO8601, prev.End)
	if err != nil {
		r.logger.Errorf("Error parsing previous node end time = \"%s\"", prev.End)
	}

	elapsed := end.Sub(start)
	waiting := start.Sub(prevEnd)

	tags := map[string]string{
		"workflow": r.cfg.WorkflowName,
		"version":  r.cfg.Version,
		"node":     r.cfg.NodeName,
	}

	fields := map[string]interface{}{
		"tracking_id": reqMsg.TrackingId,
		"node_from":   prev.NodeName,
		"elapsed_ms":  elapsed.Seconds() * 1000,
		"waiting_ms":  waiting.Seconds() * 1000,
	}

	r.handlerContext.Measurement.Save("node_elapsed_time", fields, tags)

	if isExitpoint {
		entrypoint := reqMsg.Tracking[0]
		entrypointStart, err := time.Parse(ISO8601, entrypoint.Start)
		if err != nil {
			r.logger.Errorf("Error parsing entrypoint start time = \"%s\"", entrypoint.Start)
		}
		elapsed = end.Sub(entrypointStart)

		tags = map[string]string{
			"workflow": r.cfg.WorkflowName,
			"version":  r.cfg.Version,
		}

		fields = map[string]interface{}{
			"tracking_id": reqMsg.TrackingId,
			"elapsed_ms":  elapsed.Seconds() * 1000,
		}

		r.handlerContext.Measurement.Save("workflow_elapsed_time", fields, tags)
	}
}

func sizeInKB(s []byte) string {
	return fmt.Sprintf("%.2f KB", float32(len(s))/1024)
}
