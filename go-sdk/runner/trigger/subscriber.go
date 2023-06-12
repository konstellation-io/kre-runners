package trigger

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

func (tr *TriggerRunner) startSubscriber() {
	var subscriptions []*nats.Subscription

	defer wg.Done()

	for _, subject := range viper.GetStringSlice("KRT_NATS_INPUTS") {
		consumerName := fmt.Sprintf("%s-%s", strings.ReplaceAll(subject, ".", "-"),
			strings.ReplaceAll(strings.ReplaceAll(tr.ctx.Metadata.GetProcess(), ".", "-"), " ", "-"))

		tr.ctx.Logger.V(1).Info("Subscribing to subject",
			"Subject", subject, "Queue group", consumerName)

		s, err := tr.jetstream.QueueSubscribe(
			subject,
			consumerName,
			tr.processMessage,
			nats.DeliverNew(),
			nats.Durable(consumerName),
			nats.ManualAck(),
			nats.AckWait(22*time.Hour),
		)
		if err != nil {
			tr.ctx.Logger.Error(err, "Error subscribing to NATS subject", "Subject", subject)
			os.Exit(1)
		}
		subscriptions = append(subscriptions, s)
		tr.ctx.Logger.V(1).Info("Listening to subject", "Subject", subject, "Queue group", consumerName)
	}

	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	// Handle shutdown
	tr.ctx.Logger.Info("Shutdown signal received")
	for _, s := range subscriptions {
		err := s.Unsubscribe()
		if err != nil {
			tr.ctx.Logger.Error(err, "Error unsubscribing from the NATS subject", "Subject", s.Subject)
			os.Exit(1)
		}
	}
}

func (tr *TriggerRunner) processMessage(msg *nats.Msg) {
	var (
		start = time.Now().UTC()
	)

	tr.ctx.Logger.V(1).Info("New message received")

	requestMsg, err := tr.newRequestMessage(msg.Data)
	if err != nil {
		errMsg := fmt.Sprintf("Error parsing msg.data coming from subject %s because is not a valid protobuf: %s", msg.Subject, err)
		tr.processRunnerError(msg, errMsg, requestMsg.RequestId, start, requestMsg.FromNode)
		return
	}

	tr.ctx.Logger.Info("New message received",
		"Subject", msg.Subject, "Request ID", requestMsg.RequestId)

	// Make a shallow copy of the ctx object to set inside the request msg.
	hCtx := context.ShallowCopyWithRequest(tr.ctx, requestMsg)

	if tr.responseHandler == nil {
		errMsg := "Error missing handler"
		tr.processRunnerError(msg, errMsg, requestMsg.RequestId, start, requestMsg.FromNode)
		return
	}

	err = tr.responseHandler(hCtx, requestMsg.Payload)
	if err != nil {
		errMsg := fmt.Sprintf("Error in node %q executing handler for node %q: %s",
			tr.ctx.Metadata.GetProcess(), requestMsg.FromNode, err)
		tr.processRunnerError(msg, errMsg, requestMsg.RequestId, start, requestMsg.FromNode)
		return
	}

	// Tell NATS we don't need to receive the message anymore, and we are done processing it.
	ackErr := msg.Ack()
	if ackErr != nil {
		tr.ctx.Logger.Error(ackErr, errors.ErrMsgAck)
	}

	// end := time.Now().UTC() // TODO add metrics
	// tr.saveElapsedTime(start, end, requestMsg.FromNode, true) //TODO add metrics
}

func (tr *TriggerRunner) processRunnerError(msg *nats.Msg, errMsg string, requestID string, start time.Time, fromNode string) {
	ackErr := msg.Ack()
	if ackErr != nil {
		tr.ctx.Logger.Error(ackErr, errors.ErrMsgAck)
	}

	tr.ctx.Logger.V(1).Info(errMsg)
	tr.publishError(requestID, errMsg)

	// end := time.Now().UTC() // TODO add metrics
	//tr.saveElapsedTime(start, end, fromNode, false) //TODO add metrics
}

func (tr *TriggerRunner) newRequestMessage(data []byte) (*kai.KaiNatsMessage, error) {
	requestMsg := &kai.KaiNatsMessage{}

	var err error
	if common.IsCompressed(data) {
		data, err = common.UncompressData(data)
		if err != nil {
			tr.ctx.Logger.Error(err, "Error reading compressed message")
			return nil, err
		}
	}

	err = proto.Unmarshal(data, requestMsg)

	return requestMsg, err
}

func (tr *TriggerRunner) publishError(requestID, errMsg string) {
	responseMsg := &kai.KaiNatsMessage{
		RequestId:   requestID,
		Error:       errMsg,
		FromNode:    viper.GetString("KRT_NODE_NAME"),
		MessageType: kai.MessageType_ERROR,
	}
	tr.publishResponse(responseMsg, "")
}

func (tr *TriggerRunner) publishResponse(responseMsg *kai.KaiNatsMessage, channel string) {
	outputSubject := tr.getOutputSubject(channel)

	outputMsg, err := proto.Marshal(responseMsg)
	if err != nil {
		tr.ctx.Logger.Error(err, "Error generating output result because handler result is not a serializable Protobuf")
		return
	}

	outputMsg, err = tr.prepareOutputMessage(outputMsg)
	if err != nil {
		tr.ctx.Logger.Error(err, "Error preparing output message")
		return
	}

	tr.ctx.Logger.V(1).Info("Publishing response", "Subject", outputSubject)

	_, err = tr.jetstream.Publish(outputSubject, outputMsg)
	if err != nil {
		tr.ctx.Logger.Error(err, "Error publishing output")
	}
}

func (tr *TriggerRunner) getOutputSubject(channel string) string {
	outputSubject := viper.GetString("KRT_NATS_OUTPUT")
	if channel != "" {
		return fmt.Sprintf("%s.%s", outputSubject, channel)
	}
	return outputSubject
}

// prepareOutputMessage will check the length of the message and compress it if necessary.
// Fails on compressed messages bigger than the threshold.
func (tr *TriggerRunner) prepareOutputMessage(msg []byte) ([]byte, error) {
	maxSize, err := tr.getMaxMessageSize()
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
		tr.ctx.Logger.V(1).Info("Compressed message exceeds maximum size allowed",
			"Current message size", sizeInMB(lenOutMsg),
			"Compressed message size", sizeInMB(maxSize))
		return nil, errors.ErrMessageToBig
	}

	tr.ctx.Logger.Info("Message prepared",
		"Current message size", sizeInMB(lenOutMsg),
		"Compressed message size", sizeInMB(maxSize))

	return outMsg, nil
}

// TODO this code is duplicated in the messaging package, refactor it
func (tr *TriggerRunner) getMaxMessageSize() (int64, error) {
	streamInfo, err := tr.jetstream.StreamInfo(viper.GetString("KRT_NATS_STREAM"))
	if err != nil {
		return 0, fmt.Errorf("error getting stream's max message size: %w", err)
	}

	streamMaxSize := int64(streamInfo.Config.MaxMsgSize)
	serverMaxSize := tr.nats.MaxPayload()

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
