package kre

import "github.com/konstellation-io/kre-runners/kre-go/config"

type Manager struct {
	handlers       map[string]Handler
	defaultHandler Handler
}

type ContextManager config.Config

func NewManager(defaultHandler Handler, customHandlers map[string]Handler) *Manager {
	return &Manager{
		handlers:       customHandlers,
		defaultHandler: defaultHandler,
	}
}

func (s *Manager) GetHandler(selector string) Handler {
	h, ok := s.handlers[selector]
	if !ok {
		return s.defaultHandler
	}
	return h
}
