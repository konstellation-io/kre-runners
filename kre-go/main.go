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
func Start(handlerInit HandlerInit, handler Handler) {
	logger := simplelogger.New(simplelogger.LevelInfo)
	cfg := config.NewConfig(logger)

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
	runner := NewRunner(logger, cfg, nc, js, handler, handlerInit, mongoM)
	logger.Infof("Listening to '%s' subject...", cfg.NATS.InputSubject)

	subscriptions, err := handleSubscriptions(js, cfg, runner)
	if err != nil {
		logger.Errorf("Error subscribing to NATS subjects: %s", err)
		os.Exit(1)
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

func handleSubscriptions(js nats.JetStreamContext, cfg config.Config, runner *Runner) ([]*nats.Subscription, error) {
	var (
		subscribeTo   []string
		subscriptions []*nats.Subscription
	)

	if len(cfg.NATS.InputSubjects) > 0 {
		subscribeTo = cfg.NATS.InputSubjects
	} else { //fallback to v1
		subscribeTo = append(subscribeTo, cfg.NATS.InputSubject)
	}

	for _, subject := range subscribeTo {
		s, err := js.QueueSubscribe(subject, cfg.NodeName, runner.ProcessMessage, nats.DeliverNew(), nats.Durable(cfg.NodeName), nats.ManualAck())
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, s)
	}

	return subscriptions, nil
}
