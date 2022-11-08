package config

import (
	"os"
	"strings"

	"github.com/konstellation-io/kre/libs/simplelogger"
)

type Config struct {
	WorkflowName string
	RuntimeID    string
	VersionID    string
	Version      string
	NodeName     string
	BasePath     string
	IsExitpoint  bool
	NATS         ConfigNATS
	MongoDB      MongoDB
	InfluxDB     InfluxDB
	Handlers     Handlers
}

type MongoDB struct {
	Address     string
	DataDBName  string
	ConnTimeout int
}

type ConfigNATS struct {
	Server             string
	Stream             string
	InputSubjects      []string
	OutputSubject      string
	MongoWriterSubject string
	ExitpointSubject   string
}

type InfluxDB struct {
	URI string
}

type Handlers struct {
	DefaultHandlerKey string
}

func NewConfig(logger *simplelogger.SimpleLogger) Config {
	return Config{
		WorkflowName: getCfgFromEnv(logger, "KRT_WORKFLOW_NAME"),
		RuntimeID:    getCfgFromEnv(logger, "KRT_RUNTIME_ID"),
		VersionID:    getCfgFromEnv(logger, "KRT_VERSION_ID"),
		Version:      getCfgFromEnv(logger, "KRT_VERSION"),
		NodeName:     getCfgFromEnv(logger, "KRT_NODE_NAME"),
		BasePath:     getCfgFromEnv(logger, "KRT_BASE_PATH"),
		IsExitpoint:  getCfgBoolFromEnv(logger, "KRT_IS_EXITPOINT"),
		NATS: ConfigNATS{
			Server:             getCfgFromEnv(logger, "KRT_NATS_SERVER"),
			Stream:             getCfgFromEnv(logger, "KRT_NATS_STREAM"),
			InputSubjects:      getSubscriptionsFromEnv(logger, "KRT_NATS_INPUTS"),
			OutputSubject:      getCfgFromEnv(logger, "KRT_NATS_OUTPUT"),
			MongoWriterSubject: getCfgFromEnv(logger, "KRT_NATS_MONGO_WRITER"),
			ExitpointSubject:   getCfgFromEnv(logger, "KRT_NATS_EXITPOINT_SUBJECT"),
		},
		MongoDB: MongoDB{
			Address:     getCfgFromEnv(logger, "KRT_MONGO_URI"),
			DataDBName:  "data",
			ConnTimeout: 120,
		},
		InfluxDB: InfluxDB{
			URI: getCfgFromEnv(logger, "KRT_INFLUX_URI"),
		},
		Handlers: Handlers{
			DefaultHandlerKey: "default",
		},
	}
}

func getCfgFromEnv(logger *simplelogger.SimpleLogger, name string) string {
	val, ok := os.LookupEnv(name)
	if !ok {
		logger.Errorf("Error reading config: the '%s' env var is missing", name)
		os.Exit(1)
	}
	return val
}

func getCfgBoolFromEnv(logger *simplelogger.SimpleLogger, name string) bool {
	val := getCfgFromEnv(logger, name)
	if val == "true" {
		return true
	}
	return false
}

// getSubscriptionsFromEnv will take the env variable and split it by ","
func getSubscriptionsFromEnv(logger *simplelogger.SimpleLogger, name string) []string {
	val, ok := os.LookupEnv(name)
	if !ok {
		logger.Errorf("Error reading config: the '%s' env var is missing", name)
		os.Exit(1)
	}

	return strings.Split(val, ",")
}
