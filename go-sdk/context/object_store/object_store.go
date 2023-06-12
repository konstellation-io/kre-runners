package object_store

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/internal/errors"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

type ObjectStore struct {
	logger          logr.Logger
	objStore        nats.ObjectStore
	objectStoreName string
}

func NewObjectStore(logger logr.Logger, jetstream nats.JetStreamContext) (*ObjectStore, error) {
	objectStoreName := viper.GetString("krt_nats_object_store")
	objStore, err := initObjectStoreDeps(logger, jetstream, objectStoreName)
	if err != nil {
		return nil, err
	}

	return &ObjectStore{
		logger:          logger,
		objStore:        objStore,
		objectStoreName: objectStoreName,
	}, nil
}

func initObjectStoreDeps(logger logr.Logger, jetstream nats.JetStreamContext, objectStoreName string) (nats.ObjectStore, error) {
	if objectStoreName != "" {
		objStore, err := jetstream.ObjectStore(objectStoreName)
		if err != nil {
			return nil, fmt.Errorf("error initializing object store: %w", err)
		}

		return objStore, nil

	} else {
		logger.Info("Object store not defined. Skipping object store initialization.")
		return nil, nil
	}
}

func (ob ObjectStore) Get(key string) ([]byte, error) {
	if ob.objStore == nil {
		return nil, errors.ErrUndefinedObjectStore
	}

	response, err := ob.objStore.GetBytes(key)
	if err != nil {
		return nil, fmt.Errorf("error retrieving object with key %s from the object store: %w", key, err)
	}

	ob.logger.V(1).Info("File successfully retrieved from object store", key, ob.objectStoreName)

	return response, nil
}

func (ob ObjectStore) Save(key string, payload []byte) error {
	if ob.objStore == nil {
		return errors.ErrUndefinedObjectStore
	}
	if payload == nil {
		return errors.ErrEmptyPayload
	}

	_, err := ob.objStore.PutBytes(key, payload)
	if err != nil {
		return fmt.Errorf("error storing object to the object store: %w", err)
	}

	ob.logger.V(1).Info("File successfully stored in object store", key, ob.objectStoreName)

	return nil
}

func (ob ObjectStore) Delete(key string) error {
	if ob.objStore == nil {
		return errors.ErrUndefinedObjectStore
	}

	err := ob.objStore.Delete(key)
	if err != nil {
		return fmt.Errorf("error retrieving object with key %s from the object store: %w", key, err)
	}

	ob.logger.V(1).Info("File successfully deleted in object store", key, ob.objectStoreName)

	return nil
}
