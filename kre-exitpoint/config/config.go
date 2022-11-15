package config

import (
	"os"

	"github.com/konstellation-io/kre/libs/simplelogger"
)

type Config struct {
	LastNodeName string
}

func NewConfig(logger *simplelogger.SimpleLogger) Config {
	return Config{
		LastNodeName: getCfgFromEnv(logger, "KRT_LAST_NODE_NAME"),
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
