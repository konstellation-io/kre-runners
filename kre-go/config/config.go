package config

import (
	"encoding/json"
	"os"

	"github.com/konstellation-io/kre/libs/simplelogger"
)

type Config struct {
	WorkflowName string
	RuntimeID    string
	VersionID    string
	Version      string
	NodeName     string
	BasePath     string
	NATS         ConfigNATS
	MongoDB      MongoDB
	InfluxDB     InfluxDB
	IsLastNode   bool
}

type MongoDB struct {
	Address     string
	DataDBName  string
	ConnTimeout int
}

type ConfigNATS struct {
	Server             string
	Stream             string
	InputSubject       string
	InputSubjects      []string
	OutputSubject      string
	MongoWriterSubject string
	EntrypointSubject  string
}

type InfluxDB struct {
	URI string
}

func NewConfig(logger *simplelogger.SimpleLogger) Config {
	return Config{
		WorkflowName: getCfgFromEnv(logger, "KRT_WORKFLOW_NAME"),
		RuntimeID:    getCfgFromEnv(logger, "KRT_RUNTIME_ID"),
		VersionID:    getCfgFromEnv(logger, "KRT_VERSION_ID"),
		Version:      getCfgFromEnv(logger, "KRT_VERSION"),
		NodeName:     getCfgFromEnv(logger, "KRT_NODE_NAME"),
		BasePath:     getCfgFromEnv(logger, "KRT_BASE_PATH"),
		IsLastNode:   getCfgBoolFromEnv(logger, "KRT_IS_LAST_NODE"),
		NATS: ConfigNATS{
			Server:             getCfgFromEnv(logger, "KRT_NATS_SERVER"),
			Stream:             getCfgFromEnv(logger, "KRT_NATS_STREAM"),
			InputSubject:       getCfgFromEnv(logger, "KRT_NATS_INPUT"),
			InputSubjects:      getSubscriptionsFromEnv(logger, "KRT_NATS_INPUTS"),
			OutputSubject:      getCfgFromEnv(logger, "KRT_NATS_OUTPUT"),
			MongoWriterSubject: getCfgFromEnv(logger, "KRT_NATS_MONGO_WRITER"),
			EntrypointSubject:  getCfgFromEnv(logger, "KRT_NATS_ENTRYPOINT_SUBJECT"),
		},
		MongoDB: MongoDB{
			Address:     getCfgFromEnv(logger, "KRT_MONGO_URI"),
			DataDBName:  "data",
			ConnTimeout: 120,
		},
		InfluxDB: InfluxDB{
			URI: getCfgFromEnv(logger, "KRT_INFLUX_URI"),
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

func getSubscriptionsFromEnv(logger *simplelogger.SimpleLogger, name string) []string {
	subscriptions := make([]string, 0)
	val, ok := os.LookupEnv(name)
	if ok {
		if err := json.Unmarshal([]byte(val), &subscriptions); err != nil {
			logger.Errorf("Error reading config: cannot unmarshal '%s' env var to array of strings", name)
			os.Exit(1)
		}
	}
	return subscriptions
}
