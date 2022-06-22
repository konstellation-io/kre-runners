package kre

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/konstellation-io/kre-runners/kre-go/config"
	"github.com/konstellation-io/kre-runners/kre-go/mongodb"
	"github.com/konstellation-io/kre/libs/simplelogger"
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
	handler        Handler
	handlerContext *HandlerContext
}

// NewRunner creates a new Runner instance.
func NewRunner(logger *simplelogger.SimpleLogger, cfg config.Config, nc *nats.Conn, handler Handler, handlerInit HandlerInit, mongoM *mongodb.MongoDB) *Runner {
	runner := &Runner{
		logger:  logger,
		cfg:     cfg,
		nc:      nc,
		handler: handler,
	}

	// Create handler context
	c := NewHandlerContext(cfg, nc, mongoM, logger, runner.earlyReply)
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

	// Validations
	if requestMsg.Reply == "" && msg.Reply == "" {
		r.logger.Error("the reply subject was not found")
		return
	}

	// If the "msg.Reply" has a value, it means that is the first message of the workflow.
	// In other words, it comes from the initial entrypoint request.
	// In this case we set this value into the "requestMsg.Reply" field in order
	// to be propagated for the rest of workflow nodes.
	if msg.Reply != "" {
		requestMsg.Reply = msg.Reply
	}

	r.logger.Infof("Received a message on '%s' with final reply '%s'", msg.Subject, requestMsg.Reply)

	// Make a shallow copy of the ctx object to set inside the request msg.
	hCtx := r.handlerContext
	hCtx.reqMsg = requestMsg

	// Execute the handler function sending context and the payload.
	handlerResult, err := r.handler(hCtx, requestMsg.Payload)
	if err != nil {
		r.stopWorkflowReturningErr(err, requestMsg.Reply)
		return
	}

	end := time.Now().UTC()

	// Save the elapsed time for this node and for the workflow if it is the last node.
	isLastNode := r.cfg.NATS.OutputSubject == ""
	r.saveElapsedTime(requestMsg, start, end, isLastNode)

	// Ignore send reply if the msg was replied previously.
	if isLastNode && requestMsg.Replied {
		if handlerResult != nil {
			r.logger.Info("ignoring the last node response because the message was replied previously")
		}

		return
	}

	// Generate a KreNatsMessage response.
	responseMsg, err := r.newResponseMsg(handlerResult, requestMsg, start, end, requestMsg.Reply)
	if err != nil {
		r.stopWorkflowReturningErr(err, requestMsg.Reply)
		return
	}

	// Publish the response message to the output subject.
	outputSubject := r.getOutputSubject(requestMsg.Reply, isLastNode)
	r.publishResponse(outputSubject, responseMsg)
}

// stopWorkflowReturningErr publishes a error message to the final reply subject
// in order to stop the workflow execution. So the next nodes will be ignored and the
// gRPC response will be an exception.
func (r *Runner) stopWorkflowReturningErr(err error, replySubject string) {
	r.logger.Errorf("Error executing handler: %s", err)

	errMsg := &KreNatsMessage{
		Error: fmt.Sprintf("error in '%s': %s", r.cfg.NodeName, err),
	}
	replyErrMsg, err := proto.Marshal(errMsg)
	if err != nil {
		r.logger.Errorf("Error generating error output because it is not a serializable Protobuf: %s", err)
		return
	}

	err = r.nc.Publish(replySubject, replyErrMsg)
	if err != nil {
		r.logger.Errorf("Error publishing output: %s", err)
	}
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
func (r *Runner) newResponseMsg(handlerResult proto.Message, requestMsg *KreNatsMessage, start time.Time, end time.Time, replySubject string) (*KreNatsMessage, error) {
	payload, err := ptypes.MarshalAny(handlerResult)
	if err != nil {
		return nil, fmt.Errorf("the handler result is not a valid protobuf: %w", err)
	}

	tracking := append(requestMsg.Tracking, &KreNatsMessage_Tracking{
		NodeName: r.cfg.NodeName,
		Start:    start.Format(ISO8601),
		End:      end.Format(ISO8601),
	})

	responseMsg := &KreNatsMessage{
		Replied:    requestMsg.Replied,
		Reply:      replySubject,
		TrackingId: requestMsg.TrackingId,
		Tracking:   tracking,
		Payload:    payload,
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

// getOutputSubject gets the output subject depending if this node is the last one.
func (r *Runner) getOutputSubject(replySubject string, isLastNode bool) string {
	var outputSubject string

	if isLastNode {
		outputSubject = replySubject
	} else {
		outputSubject = r.cfg.NATS.OutputSubject
	}

	return outputSubject
}

// publishResponse publishes the response in the NATS output subject.
func (r *Runner) publishResponse(outputSubject string, responseMsg *KreNatsMessage) {
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

	err = r.nc.Publish(outputSubject, outputMsg)
	if err != nil {
		r.logger.Errorf("Error publishing output: %s", err)
	}
}

func (r Runner) earlyReply(subject string, response proto.Message) error {
	payload, err := ptypes.MarshalAny(response)
	if err != nil {
		return fmt.Errorf("the handler result is not a valid protobuf: %w", err)
	}

	res := &KreNatsMessage{
		Payload: payload,
	}

	r.publishResponse(subject, res)

	return nil
}

// saveElapsedTime stores in InfluxDB the elapsed time for the current node and the total elapsed time of the
// complete workflow if it is the last node.
func (r *Runner) saveElapsedTime(reqMsg *KreNatsMessage, start time.Time, end time.Time, isLastNode bool) {
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

	if isLastNode {
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
