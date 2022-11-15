package kre

type HandlerManager struct {
	handlers       map[string]Handler
	defaultHandler Handler
}

func NewHandlerManager(defaultHandler Handler, customHandlers map[string]Handler) *HandlerManager {
	return &HandlerManager{
		handlers:       customHandlers,
		defaultHandler: defaultHandler,
	}
}

func (s *HandlerManager) GetHandler(selector string) Handler {
	h, ok := s.handlers[selector]
	if !ok {
		return s.defaultHandler
	}
	return h
}
