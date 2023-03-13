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
	NATS         ConfigNATS
	MongoDB      MongoDB
	InfluxDB     InfluxDB
}

type MongoDB struct {
	Address     string
	DataDBName  string
	ConnTimeout int
}

type ConfigNATS struct {
	Server                    string
	Stream                    string
	InputSubjects             []string
	OutputSubject             string
	ObjectStoreName           string
	KeyValueStoreProjectName  string
	KeyValueStoreWorkflowName string
	KeyValueStoreNodeName     string
	MongoWriterSubject        string
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
		NATS: ConfigNATS{
			Server:                    getCfgFromEnv(logger, "KRT_NATS_SERVER"),
			Stream:                    getCfgFromEnv(logger, "KRT_NATS_STREAM"),
			InputSubjects:             getSubscriptionsFromEnv(logger, "KRT_NATS_INPUTS"),
			OutputSubject:             getCfgFromEnv(logger, "KRT_NATS_OUTPUT"),
			ObjectStoreName:           getOptCfgFromEnv(logger, "KRT_NATS_OBJECT_STORE"),
			KeyValueStoreProjectName:  getCfgFromEnv(logger, "KRT_NATS_KEY_VALUE_STORE_PROJECT"),
			KeyValueStoreWorkflowName: getCfgFromEnv(logger, "KRT_NATS_KEY_VALUE_STORE_WORKFLOW"),
			KeyValueStoreNodeName:     getCfgFromEnv(logger, "KRT_NATS_KEY_VALUE_STORE_NODE"),
			MongoWriterSubject:        getCfgFromEnv(logger, "KRT_NATS_MONGO_WRITER"),
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
	val := getOptCfgFromEnv(logger, name)
	if val == "" {
		logger.Errorf("Error reading config: the %q env var is missing", name)
		os.Exit(1)
	}
	return val
}

func getOptCfgFromEnv(logger *simplelogger.SimpleLogger, name string) string {
	val, ok := os.LookupEnv(name)
	if !ok {
		logger.Infof("The %q env var is missing", name)
		return ""
	}
	return val
}

func getSubscriptionsFromEnv(logger *simplelogger.SimpleLogger, name string) []string {
	val, ok := os.LookupEnv(name)
	if !ok {
		logger.Errorf("Error reading config: the %q env var is missing", name)
		os.Exit(1)
	}

	return strings.Split(val, ",")
}
