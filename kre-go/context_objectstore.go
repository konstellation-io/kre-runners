package kre

import (
	"fmt"

	"github.com/nats-io/nats.go"

	"github.com/konstellation-io/kre-runners/kre-go/v4/config"
	"github.com/konstellation-io/kre-runners/kre-go/v4/internal/errors"
	"github.com/konstellation-io/kre/libs/simplelogger"
)

type contextObjectStore struct {
	cfg      config.Config
	logger   *simplelogger.SimpleLogger
	objStore nats.ObjectStore
}

func NewContextObjectStore(
	cfg config.Config,
	logger *simplelogger.SimpleLogger,
	objStore nats.ObjectStore,
) *contextObjectStore {
	return &contextObjectStore{
		cfg:      cfg,
		logger:   logger,
		objStore: objStore,
	}
}

// Save stores the given payload in the Object Store with the given key as identifier
func (c *contextObjectStore) Save(key string, payload []byte) error {
	if c.objStore == nil {
		return errors.ErrUndefinedObjectStore
	}
	if payload == nil {
		return errors.ErrEmptyPayload
	}

	_, err := c.objStore.PutBytes(key, payload)
	if err != nil {
		return fmt.Errorf("error storing object to the object store: %w", err)
	}

	c.logger.Debugf("File with key %q successfully stored in object store %q", key, c.cfg.NATS.ObjectStoreName)

	return nil
}

// Get retrieves the object stored in the node's object store
func (c *contextObjectStore) Get(key string) ([]byte, error) {
	if c.objStore == nil {
		return nil, errors.ErrUndefinedObjectStore
	}

	response, err := c.objStore.GetBytes(key)
	if err != nil {
		return nil, fmt.Errorf("error retrieving object with key %s from the object store: %w", key, err)
	}

	c.logger.Debugf("File with key %q successfully retrieved from object store %q", key, c.cfg.NATS.ObjectStoreName)

	return response, nil
}

// Delete removes the object stored in the node's object store
func (c *contextObjectStore) Delete(key string) error {
	if c.objStore == nil {
		return errors.ErrUndefinedObjectStore
	}

	err := c.objStore.Delete(key)
	if err != nil {
		return fmt.Errorf("error retrieving object with key %s from the object store: %w", key, err)
	}

	c.logger.Debugf("File with key %q successfully deleted in object store %q", key, c.cfg.NATS.ObjectStoreName)

	return nil
}
