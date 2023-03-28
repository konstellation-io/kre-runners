package kre

import (
	"fmt"

	"github.com/konstellation-io/kre-runners/kre-go/v4/config"
	"github.com/konstellation-io/kre/libs/simplelogger"
	"github.com/nats-io/nats.go"
)

type Scope string

const (
	ScopeProject  Scope = "project"
	ScopeWorkflow Scope = "workflow"
	ScopeNode     Scope = "node"
)

type contextConfiguration struct {
	kvStoresMap map[Scope]nats.KeyValue
}

func NewContextConfiguration(
	cfg config.Config,
	logger *simplelogger.SimpleLogger,
	js nats.JetStreamContext,
) (*contextConfiguration, error) {
	kvStoresMap, err := initKVStoresMap(cfg, logger, js)
	if err != nil {
		return nil, err
	}

	return &contextConfiguration{
		kvStoresMap: kvStoresMap,
	}, nil
}

func initKVStoresMap(
	cfg config.Config,
	logger *simplelogger.SimpleLogger,
	js nats.JetStreamContext,
) (map[Scope]nats.KeyValue, error) {
	kvStoresMap := make(map[Scope]nats.KeyValue, 3)

	kvStore, err := js.KeyValue(cfg.NATS.KeyValueStoreProjectName)
	if err != nil {
		return nil, err
	}
	kvStoresMap[ScopeProject] = kvStore

	kvStore, err = js.KeyValue(cfg.NATS.KeyValueStoreWorkflowName)
	if err != nil {
		return nil, err
	}
	kvStoresMap[ScopeWorkflow] = kvStore

	kvStore, err = js.KeyValue(cfg.NATS.KeyValueStoreNodeName)
	if err != nil {
		return nil, err
	}
	kvStoresMap[ScopeNode] = kvStore

	return kvStoresMap, nil
}

// Set set the given key and value to an optional scoped key-value storage,
// or the default key-value storage (Node's) if not given any.
func (cc *contextConfiguration) Set(key, value string, scopeOpt ...Scope) error {
	scope := cc.getOptionalScope(scopeOpt, ScopeNode)

	kvStore, ok := cc.kvStoresMap[scope]
	if !ok {
		return fmt.Errorf("could not find key value store given scope %q", scope)
	}

	_, err := kvStore.PutString(key, value)
	if err != nil {
		return fmt.Errorf("error storing value with key %q to the key-value store: %w", key, err)
	}

	return nil
}

// Get retrieves the configuration given a key from an optional scoped key-value storage,
// if no scoped key-value storage is given it will search in all the scopes starting by Node then upwards.
func (cc *contextConfiguration) Get(key string, scopeOpt ...Scope) (string, error) {
	if len(scopeOpt) > 0 {
		return cc.getConfigFromScope(key, scopeOpt[0])
	} else {
		allScopesInOrder := []Scope{ScopeNode, ScopeWorkflow, ScopeProject}
		for _, scope := range allScopesInOrder {
			config, err := cc.getConfigFromScope(key, scope)

			if err == nil {
				return config, nil
			}
		}

		return "", fmt.Errorf("error retrieving config with key %q, not found in any key-value store", key)
	}
}

func (cc *contextConfiguration) getConfigFromScope(key string, scope Scope) (string, error) {
	value, err := cc.kvStoresMap[scope].Get(key)

	if err != nil {
		return "", fmt.Errorf("error retrieving config with key %q from the key-value store: %w", key, err)
	}

	return string(value.Value()), nil
}

// Delete retrieves the configuration given a key from an optional scoped key-value storage,
// if no key-value storage is given it will use the default one (Node's).
func (cc *contextConfiguration) Delete(key string, scopeOpt ...Scope) error {
	scope := cc.getOptionalScope(scopeOpt, ScopeNode)

	kvStore, ok := cc.kvStoresMap[scope]
	if !ok {
		return fmt.Errorf("could not find key value store given scope %q", scope)
	}

	err := kvStore.Delete(key)
	if err != nil {
		return fmt.Errorf("error deleting value with key %q from the key-value store: %w", key, err)
	}

	return nil
}

func (cc *contextConfiguration) getOptionalScope(scopes []Scope, defaultScope Scope) Scope {
	if len(scopes) > 0 {
		return scopes[0]
	}
	return defaultScope
}
