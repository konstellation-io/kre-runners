package kre

import (
	"fmt"
	regexp2 "regexp"

	"github.com/nats-io/nats.go"

	"github.com/konstellation-io/kre-runners/kre-go/v4/config"
	"github.com/konstellation-io/kre-runners/kre-go/v4/internal/errors"
	"github.com/konstellation-io/kre/libs/simplelogger"
)

type ContextObjectStore interface {
	Save(key string, payload []byte) error
	Get(key string) ([]byte, error)
	Purge(regexp ...string) error
	List(regexp ...string) ([]string, error)
	Delete(key string) error
}

type contextObjectStore struct {
	cfg      config.Config
	logger   *simplelogger.SimpleLogger
	objStore nats.ObjectStore
}

func NewContextObjectStore(
	cfg config.Config,
	logger *simplelogger.SimpleLogger,
	js nats.JetStreamContext,
) (*contextObjectStore, error) {

	objStore, err := initObjectStore(cfg, logger, js)
	if err != nil {
		return nil, err
	}

	return &contextObjectStore{
		cfg:      cfg,
		logger:   logger,
		objStore: objStore,
	}, nil
}

func initObjectStore(cfg config.Config,
	logger *simplelogger.SimpleLogger,
	js nats.JetStreamContext,
) (nats.ObjectStore, error) {
	var objStore nats.ObjectStore
	var err error

	// Connect to ObjectStore (optional)
	if cfg.NATS.ObjectStoreName != "" {
		objStore, err = js.ObjectStore(cfg.NATS.ObjectStoreName)
		if err != nil {
			return nil, fmt.Errorf("error initializing object store: %w", err)
		}
		return objStore, nil

	} else {
		logger.Info("Object store not defined. Skipping object store initialization.")
		return nil, nil
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

func (c *contextObjectStore) Purge(regexp ...string) error {
	if c.objStore == nil {
		return errors.ErrUndefinedObjectStore
	}

	var pattern *regexp2.Regexp

	if len(regexp) > 0 && regexp[0] != "" {
		pat, err := regexp2.Compile(regexp[0])
		if err != nil {
			return fmt.Errorf("error compiling regexp: %w", err)
		}

		pattern = pat
	}

	objects, err := c.List()
	if err != nil {
		return fmt.Errorf("error listing objects from the object store: %w", err)
	}

	for _, objectName := range objects {
		if pattern == nil || pattern.MatchString(objectName) {
			err := c.objStore.Delete(objectName)

			c.logger.Debugf("Deleting object %q", objectName)
			if err != nil {
				return fmt.Errorf("error purging objects from the object store: %w", err)
			}
		}
	}

	c.logger.Debugf("Files successfully purged from object store %q", c.cfg.NATS.ObjectStoreName)

	return nil
}

func (c *contextObjectStore) List(regexp ...string) ([]string, error) {
	if c.objStore == nil {
		return nil, errors.ErrUndefinedObjectStore
	}

	objStoreList, err := c.objStore.List()
	if err != nil {
		return nil, fmt.Errorf("error listing objects from the object store: %w", err)
	}

	var pattern *regexp2.Regexp

	if len(regexp) > 0 && regexp[0] != "" {
		pat, err := regexp2.Compile(regexp[0])
		if err != nil {
			return nil, fmt.Errorf("error compiling regexp: %w", err)
		}

		pattern = pat
	}

	response := []string{}

	for _, objName := range objStoreList {
		if pattern == nil || pattern.MatchString(objName.Name) {
			response = append(response, objName.Name)
		}
	}

	c.logger.Debugf("Files successfully listed from object store %q", c.cfg.NATS.ObjectStoreName)

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
