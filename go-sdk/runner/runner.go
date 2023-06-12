package runner

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/runner/exit"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/runner/task"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/runner/trigger"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Runner struct {
	logger    logr.Logger
	nats      *nats.Conn
	jetstream nats.JetStreamContext
}

func NewRunner() *Runner {
	initializeConfiguration()

	logger := getLogger()

	nc, err := getNatsConnection(logger)
	if err != nil {
		panic(fmt.Errorf("fatal error connecting to NATS: %w", err))
	}

	js, err := getJetStreamConnection(logger, nc)
	if err != nil {
		panic(fmt.Errorf("fatal error connecting to JetStream: %w", err))
	}

	return &Runner{
		logger:    logger,
		nats:      nc,
		jetstream: js,
	}
}

func initializeConfiguration() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	viper.AutomaticEnv()

	// Set viper default values
	viper.SetDefault("KRT_BASE_PATH", "/")

	viper.SetDefault("KRT_MONGO_DB_NAME", "data")
	viper.SetDefault("KRT_MONGO_CONN_TIMEOUT", 120)
}

func getNatsConnection(logger logr.Logger) (*nats.Conn, error) {
	nc, err := nats.Connect(viper.GetString("krt_nats_url"))
	if err != nil {
		logger.Error(err, "Error connecting to NATS")
		return nil, err
	}

	return nc, nil
}

func getJetStreamConnection(logger logr.Logger, nc *nats.Conn) (nats.JetStreamContext, error) {
	js, err := nc.JetStream()
	if err != nil {
		logger.Error(err, "Error connecting to JetStream")
		return nil, err
	}

	return js, nil
}

func getLogger() logr.Logger {
	var log logr.Logger

	zapLog, err := zap.NewDevelopment()
	//zapLog, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}
	log = zapr.NewLogger(zapLog)

	return log
}

func (rn Runner) TriggerRunner() *trigger.TriggerRunner {
	return trigger.NewTriggerRunner(rn.logger, rn.nats, rn.jetstream)
}

func (rn Runner) TaskRunner() *task.TaskRunner {
	return task.NewTaskRunner(rn.logger, rn.nats, rn.jetstream)
}

func (rn Runner) ExitRunner() *exit.ExitRunner {
	return exit.NewExitRunner(rn.logger, rn.nats, rn.jetstream)
}
