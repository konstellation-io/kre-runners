package messaging

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"

	"github.com/konstellation-io/kre/libs/simplelogger"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	"github.com/konstellation-io/kre-runners/kre-go/v4/config"
	"github.com/konstellation-io/kre-runners/kre-go/v4/internal/errors"
)

const (
	MessageThreshold = 1024 * 1024
	CompressLevel    = 9
	gzipID1          = 0x1f
	gzipID2          = 0x8b
)

type NatsMessagingClient struct {
	cfg    config.Config
	logger simplelogger.SimpleLoggerInterface
	js     nats.JetStreamContext
}

func NewNatsMessagingClient(cfg config.Config, logger simplelogger.SimpleLoggerInterface, js nats.JetStreamContext) *NatsMessagingClient {
	return &NatsMessagingClient{
		cfg:    cfg,
		logger: logger,
		js:     js,
	}
}

//func (m *NatsMessagingClient) PublishResponse(responseMsg *KreNatsMessage, channel string) {
func (m *NatsMessagingClient) PublishResponse(responseMsg proto.Message, channel string) {
	outputSubject := m.getOutputSubject(channel)

	outputMsg, err := proto.Marshal(responseMsg)
	if err != nil {
		m.logger.Errorf("Error generating output result because handler result is not a serializable Protobuf: %s", err)
		return
	}

	outputMsg, err = m.prepareOutputMessage(outputMsg)
	if err != nil {
		m.logger.Errorf("Error preparing output msg: %s", err)
		return
	}

	m.logger.Infof("Publishing response to %q subject", outputSubject)

	_, err = m.js.Publish(outputSubject, outputMsg)
	if err != nil {
		m.logger.Errorf("Error publishing output: %s", err)
	}
}
func (m *NatsMessagingClient) getOutputSubject(channel string) string {
	outputSubject := m.cfg.NATS.OutputSubject
	if channel != "" {
		return fmt.Sprintf("%s.%s", outputSubject, channel)
	}
	return outputSubject
}

// prepareOutputMessage will check the length of the message and compress it if necessary.
// Fails on compressed messages bigger than the threshold.
func (m *NatsMessagingClient) prepareOutputMessage(msg []byte) ([]byte, error) {
	if len(msg) <= MessageThreshold {
		return msg, nil
	}

	outMsg, err := m.compressData(msg)
	if err != nil {
		return nil, err
	}

	if len(outMsg) > MessageThreshold {
		return nil, errors.ErrMessageToBig
	}

	m.logger.Infof("Original message size: %s. Compressed: %s", sizeInKB(msg), sizeInKB(outMsg))

	return outMsg, nil
}

// compressData creates compressed []byte.
func (m *NatsMessagingClient) compressData(data []byte) ([]byte, error) {
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
func (m *NatsMessagingClient) uncompressData(data []byte) ([]byte, error) {
	rd := bytes.NewReader(data)
	gr, err := gzip.NewReader(rd)
	if err != nil {
		return nil, err
	}

	defer gr.Close()
	return ioutil.ReadAll(gr)
}

// isCompressed check if the input string is compressed.
func (m *NatsMessagingClient) isCompressed(data []byte) bool {
	return data[0] == gzipID1 && data[1] == gzipID2
}

func sizeInKB(s []byte) string {
	return fmt.Sprintf("%.2f KB", float32(len(s))/1024)
}
