package main

import (
	"log"
	"os"

	"nats_init/config"

	"github.com/nats-io/nats.go"
)

func main() {
	log.Default().Printf("Init NATS stream")

	cfg := config.NewConfig()

	natsConn, err := nats.Connect(cfg.Server)
	if err != nil {
		log.Fatalf("Failed to connect to nats: %v", err)
		os.Exit(1)
	}
	js, err := natsConn.JetStream()
	if err != nil {
		log.Fatalf("Failed connecting to NATS JetStream: %v", err)
		os.Exit(1)
	}

	streamCfg := &nats.StreamConfig{
		Name:        cfg.Stream,
		Description: "e2e stream",
		Subjects:    []string{cfg.Stream + ".*"},
		Retention:   nats.InterestPolicy,
	}
	_, err = js.AddStream(streamCfg)
	if err != nil {
		log.Fatalf("Failed creating Stream: %v", err)
		os.Exit(1)
	}

	log.Default().Printf("Stream %s created succesfully", cfg.Stream)
}
