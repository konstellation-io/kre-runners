package kre

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/konstellation-io/kre-runners/kre-go/v4/config"
	"github.com/konstellation-io/kre-runners/kre-go/v4/mongodb"
	"github.com/konstellation-io/kre/libs/simplelogger"
)

type QueryData map[string]interface{}

type SaveDataMsg struct {
	Coll string      `json:"coll"`
	Doc  interface{} `json:"doc"`
}

type ContextDatabase interface {
	Find(collection string, query QueryData, res interface{}) error
	Save(collection string, data interface{}) error
}

type contextDatabase struct {
	cfg    config.Config
	nc     *nats.Conn
	mongoM mongodb.Manager
	logger *simplelogger.SimpleLogger
}

func NewContextDatabase(
	cfg config.Config,
	nc *nats.Conn,
	mongoM mongodb.Manager,
	logger *simplelogger.SimpleLogger,
) ContextDatabase {
	return &contextDatabase{
		cfg:    cfg,
		nc:     nc,
		mongoM: mongoM,
		logger: logger,
	}
}

// Find data from a collection of mongoDB
func (c *contextDatabase) Find(collection string, query QueryData, res interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), getDataTimeout)
	defer cancel()

	criteria := bson.M{}
	for k, v := range query {
		criteria[k] = v
	}

	return c.mongoM.Find(ctx, collection, criteria, res)
}

// Save data inside a bson struct to a collection of your choice in mongoDB
func (c *contextDatabase) Save(collection string, data interface{}) error {
	msg, err := json.Marshal(SaveDataMsg{
		Coll: collection,
		Doc:  data,
	})
	if err != nil {
		c.logger.Infof("Error generating SaveDataMsg JSON: %s", err)
	}

	_, err = c.nc.Request(c.cfg.NATS.MongoWriterSubject, msg, saveDataTimeout)
	return err
}
