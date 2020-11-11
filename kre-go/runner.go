package kre

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/nats-io/nats.go"
	"time"

	"github.com/konstellation-io/kre-runners/kre-go/config"
	"github.com/konstellation-io/kre/libs/simplelogger"
)

const ISO8601 = "2006-01-02T15:04:05.000000"

type Runner struct {
	logger         *simplelogger.SimpleLogger
	cfg            config.Config
	nc             *nats.Conn
	handler        Handler
	handlerContext *HandlerContext
}

// NewRunner creates a new Runner instance.
func NewRunner(
	logger *simplelogger.SimpleLogger,
	cfg config.Config,
	nc *nats.Conn,
	handler Handler,
	c *HandlerContext,
) *Runner {
	return &Runner{
		logger:         logger,
		cfg:            cfg,
		nc:             nc,
		handler:        handler,
		handlerContext: c,
	}
}

// ProcessMessage parses the incoming NATS message, executes the handler function and publishes
// the handler result to the output subject.
func (r *Runner) ProcessMessage(msg *nats.Msg) {
	start := time.Now().UTC()
	r.logger.Infof("Received a message on '%s' with reply '%s'", msg.Subject, msg.Reply)

	// Parse incoming message
	requestMsg := &KreNatsMessage{}
	err := proto.Unmarshal(msg.Data, requestMsg)
	if err != nil {
		r.logger.Errorf("Error parsing msg.data because is not a valid protobuf: %s", err)
		return
	}

	// Get reply from the NATS msg or the request msg
	replySubject, err := getReplySubject(msg, requestMsg)
	if err != nil {
		r.logger.Error(err.Error())
		return
	}

	// Execute the handler function sending context and the payload
	handlerResult, err := r.handler(r.handlerContext, requestMsg.Payload)
	if err != nil {
		r.stopWorkflowReturningErr(err, replySubject)
		return
	}

	// Generate a KreNatsMessage response
	responseMsg, err := r.createResponseMsg(handlerResult, requestMsg, start, replySubject)
	if err != nil {
		r.stopWorkflowReturningErr(err, replySubject)
		return
	}

	// Publish the response message to the output subject.
	r.sendOutputMsg(replySubject, responseMsg)
}

// getReplySubject gets the final reply from the nats message or the request.
func getReplySubject(msg *nats.Msg, requestMsg *KreNatsMessage) (string, error) {
	if requestMsg.Reply == "" && msg.Reply == "" {
		return "", errors.New("the reply subject was not found")
	}

	if msg.Reply != "" {
		return msg.Reply, nil
	}

	return requestMsg.Reply, nil
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

// createResponseMsg creates a KreNatsMessage maintaining the tracking ID and adding the
// handler result and the tracking information for this node.
func (r *Runner) createResponseMsg(handlerResult proto.Message, requestMsg *KreNatsMessage, start time.Time, replySubject string) (*KreNatsMessage, error) {
	payload, err := ptypes.MarshalAny(handlerResult)
	if err != nil {
		return nil, fmt.Errorf("the handler result is not a valid protobuf: %w", err)
	}

	// Add tracking info
	end := time.Now().UTC()
	tracking := append(requestMsg.Tracking, &KreNatsMessage_Tracking{
		NodeName: r.cfg.NodeName,
		Start:    start.Format(ISO8601),
		End:      end.Format(ISO8601),
	})

	responseMsg := &KreNatsMessage{
		TrackingId: requestMsg.TrackingId,
		Tracking:   tracking,
		Reply:      replySubject,
		Payload:    payload,
	}

	return responseMsg, nil
}

// getOutputSubject gets the output subject depending if this node is the last one.
func (r *Runner) getOutputSubject(replySubject string) string {
	var outputSubject string

	isLastNode := r.cfg.NATS.OutputSubject == ""
	if isLastNode {
		outputSubject = replySubject
	} else {
		outputSubject = r.cfg.NATS.OutputSubject
	}

	return outputSubject
}

// sendOutputMsg publishes the response in the NATS output subject.
func (r *Runner) sendOutputMsg(replySubject string, responseMsg *KreNatsMessage) {
	outputMsg, err := proto.Marshal(responseMsg)
	if err != nil {
		r.logger.Errorf("Error generating output result because handler result is not a serializable Protobuf: %s", err)
		return
	}

	outputSubject := r.getOutputSubject(replySubject)
	r.logger.Infof("Publishing response to '%s' subject", outputSubject)

	err = r.nc.Publish(outputSubject, outputMsg)
	if err != nil {
		r.logger.Errorf("Error publishing output: %s", err)
	}
}
