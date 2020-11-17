package kre

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/konstellation-io/kre/libs/simplelogger"
	"github.com/nats-io/nats.go"

	"github.com/konstellation-io/kre-runners/kre-go/config"
)

const (
	ISO8601          = "2006-01-02T15:04:05.000000"
	MessageThreshold = 1024 * 1024
	CompressLevel    = 9
	gzipID1          = 0x1f
	gzipID2          = 0x8b
)

var ErrMessageToBig = errors.New("compressed message exceeds maximum size allowed of 1 MB.")

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
	requestMsg, err := r.getRequestMessage(msg.Data)
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

// prepareOutputMessage check the length of the message and compress if necessary.
// fails on compressed messages bigger than the threshold.
func (r *Runner) prepareOutputMessage(msg []byte) ([]byte, error) {
	if len(msg) <= MessageThreshold {
		return msg, nil
	}

	out_msg, err := r.compressData(msg)
	if err != nil {
		return nil, err
	}

	if len(out_msg) > MessageThreshold {
		return nil, ErrMessageToBig
	}

	r.logger.Infof("Original message size: %s. Compressed: %s", sizeInKB(msg), sizeInKB(out_msg))

	return out_msg, nil
}

// isCompressed check if the input string is compressed.
func (r *Runner) isCompressed(data []byte) bool {
	return data[0] == gzipID1 && data[1] == gzipID2
}

// compressData creates compressed []byte.
func (r *Runner) compressData(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz, err := gzip.NewWriterLevel(&b, CompressLevel)
	if err != nil {
		return nil, err
	}
	_, err = gz.Write(data)
	if err != nil {
		return nil, err
	}
	err = gz.Close()
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// uncompressData open gzip and return uncompressed []byte.
func (r *Runner) uncompressData(data []byte) ([]byte, error) {
	rd := bytes.NewReader(data)
	gr, err := gzip.NewReader(rd)
	if err != nil {
		return nil, err
	}

	defer gr.Close()
	return ioutil.ReadAll(gr)
}

// getRequestMessage creates an instance of KreNatsMessage for the input string. decompress if necessary
func (r *Runner) getRequestMessage(data []byte) (*KreNatsMessage, error) {
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

func sizeInKB(s []byte) string {
	return fmt.Sprintf("%.2f KB", float32(len(s))/1024)
}
