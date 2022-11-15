package config

import (
	"log"
	"os"
)

type Config struct {
	Server string
	Stream string
}

func NewConfig() Config {
	return Config{
		Server: getCfgFromEnv("KRT_NATS_SERVER"),
		Stream: getCfgFromEnv("KRT_NATS_STREAM"),
	}
}

func getCfgFromEnv(name string) string {
	val, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalf("Error reading config: the '%s' env var is missing", name)
		os.Exit(1)
	}
	return val
}
