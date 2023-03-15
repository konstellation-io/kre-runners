package config

import (
	"log"
	"os"
)

type Config struct {
	Server      string
	Stream      string
	ObjectStore string
	KvsConfigs  KVSConfigs
}

type KVSConfigs struct {
	KVSProject    string
	KVSWorkflow   string
	KVSEntrypoint string
	KVSExitpoint  string
	KVSNodeA      string
	KVSNodeB      string
	KVSNodeC      string
}

func NewConfig() Config {
	return Config{
		Server:      getCfgFromEnv("KRT_NATS_SERVER"),
		Stream:      getCfgFromEnv("KRT_NATS_STREAM"),
		ObjectStore: getCfgFromEnv("KRT_NATS_OBJECT_STORE"),
		KvsConfigs: KVSConfigs{
			KVSProject:    getCfgFromEnv("KRT_NATS_KEY_VALUE_STORE_PROJECT"),
			KVSWorkflow:   getCfgFromEnv("KRT_NATS_KEY_VALUE_STORE_WORKFLOW"),
			KVSEntrypoint: getCfgFromEnv("KRT_NATS_KEY_VALUE_STORE_ENTRYPOINT"),
			KVSExitpoint:  getCfgFromEnv("KRT_NATS_KEY_VALUE_STORE_EXITPOINT"),
			KVSNodeA:      getCfgFromEnv("KRT_NATS_KEY_VALUE_STORE_NODE_A"),
			KVSNodeB:      getCfgFromEnv("KRT_NATS_KEY_VALUE_STORE_NODE_B"),
			KVSNodeC:      getCfgFromEnv("KRT_NATS_KEY_VALUE_STORE_NODE_C"),
		},
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
