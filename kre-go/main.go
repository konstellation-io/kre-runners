package kre

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/konstellation-io/kre-runners/kre-go/config"
	"github.com/konstellation-io/kre-runners/kre-go/mongodb"
	"github.com/konstellation-io/kre/libs/simplelogger"
)

// HandlerInit is executed once. It is useful to initialize variables that will be constants
// between handler calls.
type HandlerInit func(ctx *HandlerContext)

// Handler is the function executed each time a message from NATS arrives. This function must return
// the protobuf for the next node or the final response if it is the last one.
type Handler func(ctx *HandlerContext, data *anypb.Any) (proto.Message, error)

// Start receives the handler init function and the handler function
// connects to NATS and MongoDB and processes all incoming messages.
func Start(handlerInit HandlerInit, handler Handler, handlersOpt ...map[string]Handler) {
	logger := simplelogger.New(simplelogger.LevelInfo)
	cfg := config.NewConfig(logger)

	if handler == nil && handlersOpt == nil {
		logger.Errorf("Error no handlers detected")
		os.Exit(1)
	}

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

	if handler != nil && handlersOpt == nil {
		handlersOpt = make([]map[string]Handler, 1)
		handlersOpt[0] = map[string]Handler{cfg.Handlers.DefaultHandlerKey: handler}
	}

	// Handle incoming messages from NATS
	runner := NewRunner(logger, cfg, nc, js, handlersOpt[0], handlerInit, mongoM)

	var subscriptions []*nats.Subscription
	for _, subject := range cfg.NATS.InputSubjects {
		s, err := js.QueueSubscribe(subject, cfg.NodeName, runner.ProcessMessage, nats.DeliverNew(), nats.Durable(cfg.NodeName), nats.ManualAck())
		if err != nil {
			logger.Errorf("Error subscribing to NATS subject %s: %s", subject, err)
			os.Exit(1)
		}
		subscriptions = append(subscriptions, s)
		logger.Infof("Listening to '%s' subject with queue %s", subject, cfg.NodeName)
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
