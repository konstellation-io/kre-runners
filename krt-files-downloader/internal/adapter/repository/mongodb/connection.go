package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectTimeout = 30 * time.Second
const disconnectTimeout = 30 * time.Second

func (m *MongoManager) Connect(address string) error {
	log.Print("MongoDB connecting...")

	client, err := mongo.NewClient(options.Client().ApplyURI(address))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	log.Print("MongoDB ping...")

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	log.Print("MongoDB connected")

	m.client = client

	return nil
}

func (m *MongoManager) Disconnect() error {
	log.Print("MongoDB disconnecting...")

	if m.client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), disconnectTimeout)
	defer cancel()

	err := m.client.Disconnect(ctx)

	if err != nil {
		return err
	}

	log.Print("Connection to MongoDB closed.")

	return nil
}
