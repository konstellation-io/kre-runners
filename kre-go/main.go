package kre

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/konstellation-io/kre-runners/kre-go/config"
	"github.com/konstellation-io/kre-runners/kre-go/mongodb"
	"github.com/konstellation-io/kre/libs/simplelogger"
	"github.com/nats-io/nats.go"
	"os"
	"os/signal"
	"syscall"
)

// HandlerInit is executed once. It is useful to initialize variables that will be constants
// between handler calls.
type HandlerInit func(ctx *HandlerContext)

// Handler is the function executed each time a message from NATS arrives. This function must return
// the protobuf for the next node or the final response if it is the last one.
type Handler func(ctx *HandlerContext, data *any.Any) (proto.Message, error)

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

	// Connect to MongoDB
	mongoM := mongodb.NewMongoManager(cfg, logger)
	err = mongoM.Connect()
	if err != nil {
		logger.Errorf("Error connecting to MongoDB: %s", err)
		os.Exit(1)
	}

	// Create handler context
	c := NewHandlerContext(cfg, nc, mongoM, logger)
	handlerInit(c)

	// Handle incoming messages from NATS
	runner := NewRunner(logger, cfg, nc, handler, c)
	logger.Infof("Listening to '%s' subject...", cfg.NATS.InputSubject)
	s, err := nc.Subscribe(cfg.NATS.InputSubject, runner.ProcessMessage)
	if err != nil {
		logger.Errorf("Error subscribing to the NATS subject: %s", err)
		os.Exit(1)
	}

	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	// Handle shutdown
	logger.Info("Shutdown signal received")
	err = s.Unsubscribe()
	if err != nil {
		logger.Errorf("Error unsubscribing from the NATS subject: %s", err)
		os.Exit(1)
	}
}
