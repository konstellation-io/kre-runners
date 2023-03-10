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

	"github.com/konstellation-io/kre-runners/kre-go/v4/config"
	"github.com/konstellation-io/kre-runners/kre-go/v4/mongodb"
)

// HandlerInit is executed once. It is useful to initialize variables that will be constants
// between handler calls.
type HandlerInit func(ctx *HandlerContext)

// Handler is the function executed each time a message from NATS arrives.
//
// Responses, if desired, must be sent through the handlerContext's sendOutput or SendAny funcs.
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

	handlerManager := NewHandlerManager(defaultHandler, customHandler)

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

	// Connect to ObjectStore (optional)
	var objStore nats.ObjectStore
	if cfg.NATS.ObjectStoreName != "" {
		objStore, err = js.ObjectStore(cfg.NATS.ObjectStoreName)
		if err != nil {
			logger.Errorf("error binding the object store: %s", err)
			os.Exit(1)
		}
	}

	kvStoresMap := initKVStoresMap(cfg, logger, js)

	// Connect to MongoDB
	mongoM := mongodb.NewMongoManager(cfg, logger)
	err = mongoM.Connect()
	if err != nil {
		logger.Errorf("Error connecting to MongoDB: %s", err)
		os.Exit(1)
	}

	// Handle incoming messages from NATS
	runner := NewRunner(logger, cfg, nc, js, objStore, kvStoresMap, handlerManager, handlerInit, mongoM)

	var subscriptions []*nats.Subscription
	for _, subject := range cfg.NATS.InputSubjects {
		consumerName := fmt.Sprintf("%s-%s", strings.ReplaceAll(subject, ".", "-"), cfg.NodeName)

		s, err := js.QueueSubscribe(
			subject,
			consumerName,
			runner.ProcessMessage,
			nats.DeliverNew(),
			nats.Durable(consumerName),
			nats.ManualAck(),
			nats.AckWait(22*time.Hour),
		)
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

func initKVStoresMap(cfg config.Config, logger *simplelogger.SimpleLogger, js nats.JetStreamContext) map[Scope]nats.KeyValue {
	kvStoresMap := make(map[Scope]nats.KeyValue, 3)

	kvStore, err := js.KeyValue(cfg.NATS.KeyValueStoreProjectName)
	if err != nil {
		logger.Errorf("error binding the key value store: %s", err)
		os.Exit(1)
	}

	kvStoresMap[ScopeProject] = kvStore

	kvStore, err = js.KeyValue(cfg.NATS.KeyValueStoreWorkflowName)
	if err != nil {
		logger.Errorf("error binding the key value store: %s", err)
		os.Exit(1)
	}

	kvStoresMap[ScopeWorkflow] = kvStore

	kvStore, err = js.KeyValue(cfg.NATS.KeyValueStoreNodeName)
	if err != nil {
		logger.Errorf("error binding the key value store: %s", err)
		os.Exit(1)
	}

	kvStoresMap[ScopeNode] = kvStore

	return kvStoresMap
}
