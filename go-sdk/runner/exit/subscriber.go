package exit

import (
	"fmt"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/context"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/internal/common"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/internal/errors"
	kai "github.com/konstellation-io/kre-runners/go-sdk/v1/protos"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/proto"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func (er *ExitRunner) startSubscriber() {
	var subscriptions []*nats.Subscription

	defer wg.Done()

	for _, subject := range viper.GetStringSlice("KRT_NATS_INPUTS") {
		consumerName := fmt.Sprintf("%s-%s", strings.ReplaceAll(subject, ".", "-"),
			strings.ReplaceAll(strings.ReplaceAll(er.ctx.Metadata.GetProcess(), ".", "-"), " ", "-"))

		er.ctx.Logger.V(1).Info("Subscribing to subject",
			"Subject", subject, "Queue group", consumerName)

		s, err := er.jetstream.QueueSubscribe(
			subject,
			consumerName,
			er.processMessage,
			nats.DeliverNew(),
			nats.Durable(consumerName),
			nats.ManualAck(),
			nats.AckWait(22*time.Hour),
		)
		if err != nil {
			er.ctx.Logger.Error(err, "Error subscribing to NATS subject", "Subject", subject)
			os.Exit(1)
		}
		subscriptions = append(subscriptions, s)
		er.ctx.Logger.V(1).Info("Listening to subject", "Subject", subject, "Queue group", consumerName)
	}

	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	// Handle shutdown
	er.ctx.Logger.Info("Shutdown signal received")
	for _, s := range subscriptions {
		err := s.Unsubscribe()
		if err != nil {
			er.ctx.Logger.Error(err, "Error unsubscribing from the NATS subject", "Subject", s.Subject)
			os.Exit(1)
		}
	}
}

func (er *ExitRunner) processMessage(msg *nats.Msg) {
	var (
		start = time.Now().UTC()
	)

	requestMsg, err := er.newRequestMessage(msg.Data)
	if err != nil {
		errMsg := fmt.Sprintf("Error parsing msg.data coming from subject %s because is not a valid protobuf: %s", msg.Subject, err)
		er.processRunnerError(msg, errMsg, requestMsg.RequestId, start, requestMsg.FromNode)
		return
	}

	er.ctx.Logger.Info("New message received",
		"Subject", msg.Subject, "Request ID", requestMsg.RequestId)

	// Make a shallow copy of the ctx object to set inside the request msg.
	hCtx := context.ShallowCopyWithRequest(er.ctx, requestMsg)

	handler := er.getResponseHandler(strings.ToLower(requestMsg.FromNode))
	if handler == nil {
		errMsg := fmt.Sprintf("Error missing handler for node %q", requestMsg.FromNode)
		er.processRunnerError(msg, errMsg, requestMsg.RequestId, start, requestMsg.FromNode)
		return
	}

	if er.preprocessor != nil {
		err := er.preprocessor(hCtx, requestMsg.Payload)
		if err != nil {
			errMsg := fmt.Sprintf("Error in node %q executing handler preprocessor for node %q: %s",
				er.ctx.Metadata.GetProcess(), requestMsg.FromNode, err)
			er.processRunnerError(msg, errMsg, requestMsg.RequestId, start, requestMsg.FromNode)
			return
		}
	}

	err = handler(hCtx, requestMsg.Payload)
	if err != nil {
		errMsg := fmt.Sprintf("Error in node %q executing handler for node %q: %s",
			er.ctx.Metadata.GetProcess(), requestMsg.FromNode, err)
		er.processRunnerError(msg, errMsg, requestMsg.RequestId, start, requestMsg.FromNode)
		return
	}

	if er.postprocessor != nil {
		err := er.postprocessor(hCtx, requestMsg.Payload)
		if err != nil {
			errMsg := fmt.Sprintf("Error in node %q executing handler postprocessor for node %q: %s",
				er.ctx.Metadata.GetProcess(), requestMsg.FromNode, err)
			er.processRunnerError(msg, errMsg, requestMsg.RequestId, start, requestMsg.FromNode)
			return
		}
	}

	// Tell NATS we don't need to receive the message anymore, and we are done processing it.
	ackErr := msg.Ack()
	if ackErr != nil {
		er.ctx.Logger.Error(ackErr, errors.ErrMsgAck)
	}

	// end := time.Now().UTC() // TODO add metrics
	// er.saveElapsedTime(start, end, requestMsg.FromNode, true) //TODO add metrics
}

func (er *ExitRunner) processRunnerError(msg *nats.Msg, errMsg string, requestID string, start time.Time, fromNode string) {
	ackErr := msg.Ack()
	if ackErr != nil {
		er.ctx.Logger.Error(ackErr, errors.ErrMsgAck)
	}

	er.ctx.Logger.V(1).Info(errMsg)
	er.publishError(requestID, errMsg)

	// end := time.Now().UTC() // TODO add metrics
	//er.saveElapsedTime(start, end, fromNode, false) //TODO add metrics
}

func (er *ExitRunner) newRequestMessage(data []byte) (*kai.KaiNatsMessage, error) {
	requestMsg := &kai.KaiNatsMessage{}

	var err error
	if common.IsCompressed(data) {
		data, err = common.UncompressData(data)
		if err != nil {
			er.ctx.Logger.Error(err, "Error reading compressed message")
			return nil, err
		}
	}

	err = proto.Unmarshal(data, requestMsg)

	return requestMsg, err
}

func (er *ExitRunner) publishError(requestID, errMsg string) {
	responseMsg := &kai.KaiNatsMessage{
		RequestId:   requestID,
		Error:       errMsg,
		FromNode:    viper.GetString("KRT_NODE_NAME"),
		MessageType: kai.MessageType_ERROR,
	}
	er.publishResponse(responseMsg, "")
}

func (er *ExitRunner) publishResponse(responseMsg *kai.KaiNatsMessage, channel string) {
	outputSubject := er.getOutputSubject(channel)

	outputMsg, err := proto.Marshal(responseMsg)
	if err != nil {
		er.ctx.Logger.Error(err, "Error generating output result because handler result is not a serializable Protobuf")
		return
	}

	outputMsg, err = er.prepareOutputMessage(outputMsg)
	if err != nil {
		er.ctx.Logger.Error(err, "Error preparing output message")
		return
	}

	er.ctx.Logger.V(1).Info("Publishing response", "Subject", outputSubject)

	_, err = er.jetstream.Publish(outputSubject, outputMsg)
	if err != nil {
		er.ctx.Logger.Error(err, "Error publishing output")
	}
}

func (er *ExitRunner) getOutputSubject(channel string) string {
	outputSubject := viper.GetString("KRT_NATS_OUTPUT")
	if channel != "" {
		return fmt.Sprintf("%s.%s", outputSubject, channel)
	}
	return outputSubject
}

// prepareOutputMessage will check the length of the message and compress it if necessary.
// Fails on compressed messages bigger than the threshold.
func (er *ExitRunner) prepareOutputMessage(msg []byte) ([]byte, error) {
	maxSize, err := er.getMaxMessageSize()
	if err != nil {
		return nil, fmt.Errorf("error getting max message size: %s", err)
	}

	lenMsg := int64(len(msg))
	if lenMsg <= maxSize {
		return msg, nil
	}

	outMsg, err := common.CompressData(msg)
	if err != nil {
		return nil, err
	}

	lenOutMsg := int64(len(outMsg))
	if lenOutMsg > maxSize {
		er.ctx.Logger.V(1).Info("Compressed message exceeds maximum size allowed",
			"Current message size", sizeInMB(lenOutMsg),
			"Compressed message size", sizeInMB(maxSize))
		return nil, errors.ErrMessageToBig
	}

	er.ctx.Logger.Info("Message prepared",
		"Current message size", sizeInMB(lenOutMsg),
		"Compressed message size", sizeInMB(maxSize))

	return outMsg, nil
}

func (er *ExitRunner) getResponseHandler(subject string) Handler {
	if responseHandler, ok := er.responseHandlers[subject]; ok {
		return responseHandler
	}

	// returns the default response handler, or nil if it doesn't exist
	return er.responseHandlers["default"]
}

// TODO this code is duplicated in the messaging package, refactor it
func (er *ExitRunner) getMaxMessageSize() (int64, error) {
	streamInfo, err := er.jetstream.StreamInfo(viper.GetString("KRT_NATS_STREAM"))
	if err != nil {
		return 0, fmt.Errorf("error getting stream's max message size: %w", err)
	}

	streamMaxSize := int64(streamInfo.Config.MaxMsgSize)
	serverMaxSize := er.nats.MaxPayload()

	if streamMaxSize == -1 {
		return serverMaxSize, nil
	}

	if streamMaxSize < serverMaxSize {
		return streamMaxSize, nil
	}

	return serverMaxSize, nil
}

func sizeInMB(size int64) string {
	return fmt.Sprintf("%.1f MB", float32(size)/1024/1024)
}
