package kre

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/konstellation-io/kre/libs/simplelogger"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/konstellation-io/kre-runners/kre-go/config"
	"github.com/konstellation-io/kre-runners/kre-go/mongodb"
)

// HandlerInit is executed once. It is useful to initialize variables that will be constants
// between handler calls.
type HandlerInit func(ctx *HandlerContext)

// Handler is the function executed each time a message from NATS arrives.
// Responses must be sent through the handlerContext's sendOutput func.
type Handler func(ctx *HandlerContext, data *anypb.Any) error

// Start receives the handler init function and the handler function
// connects to NATS and MongoDB and processes all incoming messages.
func Start(handlerInit HandlerInit, defaultHandler Handler, handlersOpt ...map[string]Handler) {
	logger := simplelogger.New(simplelogger.LevelInfo)
	cfg := config.NewConfig(logger)

	noHandlersDefined := handlersOpt == nil || len(handlersOpt) < 1
	if defaultHandler == nil && noHandlersDefined {
		logger.Errorf("No handlers detected")
		os.Exit(1)
	}

	var customHandler map[string]Handler
	if len(handlersOpt) > 0 {
		customHandler = handlersOpt[0]
	}

	handlerManager := NewManager(defaultHandler, customHandler)

	// Connect to NATS
	nc, err := nats.Connect(cfg.NATS.Server)
	if err != nil {
		logger.Errorf("Error connecting to NATS: %s", err)
		os.Exit(1)
	}
	defer nc.Close()

	// Connect to JetStream
	js, err := nc.JetStream()
	if err != nil {
		logger.Errorf("Error connecting to JetStream: %s", err)
		os.Exit(1)
	}

	// Connect to MongoDB
	mongoM := mongodb.NewMongoManager(cfg, logger)
	err = mongoM.Connect()
	if err != nil {
		logger.Errorf("Error connecting to MongoDB: %s", err)
		os.Exit(1)
	}

	// Handle incoming messages from NATS
	runner := NewRunner(logger, cfg, nc, js, handlerManager, handlerInit, mongoM)

	var subscriptions []*nats.Subscription
	for _, subject := range cfg.NATS.InputSubjects {
		consumerName := fmt.Sprintf("%s-%s", strings.ReplaceAll(subject, ".", "-"), cfg.NodeName)

		s, err := js.QueueSubscribe(subject, consumerName, runner.ProcessMessage,
			nats.DeliverNew(), nats.Durable(consumerName), nats.ManualAck(), nats.AckWait(22*time.Hour))
		if err != nil {
			logger.Errorf("Error subscribing to NATS subject %s: %s", subject, err)
			os.Exit(1)
		}
		subscriptions = append(subscriptions, s)
		logger.Infof("Listening to '%s' subject with queue group %s", subject, consumerName)
	}

	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	// Handle shutdown
	logger.Info("Shutdown signal received")
	for _, s := range subscriptions {
		err = s.Unsubscribe()
		if err != nil {
			logger.Errorf("Error unsubscribing from the NATS subject %s: %s", s.Subject, err)
			os.Exit(1)
		}
	}
}
