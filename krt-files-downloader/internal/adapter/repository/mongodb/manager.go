package mongodb

import "go.mongodb.org/mongo-driver/mongo"

type MongoManager struct {
	client *mongo.Client
}

func NewMongoManger() *MongoManager {
	return &MongoManager{}
}
