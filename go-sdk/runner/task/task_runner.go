package task

import (
	"github.com/go-logr/logr"
	"github.com/nats-io/nats.go"
	"sync"

	"github.com/konstellation-io/kre-runners/go-sdk/v1/context"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/runner/common"
)

type Preprocessor common.Handler

type Handler common.Handler

type Postprocessor common.Handler

type TaskRunner struct {
	ctx              context.KaiContext
	nats             *nats.Conn
	jetstream        nats.JetStreamContext
	responseHandlers map[string]Handler
	initializer      common.Initializer
	preprocessor     Preprocessor
	handler          Handler
	postprocessor    Postprocessor
	finalizer        common.Finalizer
}

var wg sync.WaitGroup

func NewTaskRunner(logger logr.Logger, nats *nats.Conn, jetstream nats.JetStreamContext) *TaskRunner {
	return &TaskRunner{
		ctx:              context.NewContext(logger.WithName("[Task]"), nats, jetstream),
		nats:             nats,
		jetstream:        jetstream,
		responseHandlers: make(map[string]Handler),
	}
}

func (tr *TaskRunner) WithInitializer(initializer common.Initializer) *TaskRunner {
	tr.initializer = composeInitializer(initializer)
	return tr
}

func (tr *TaskRunner) WithPreprocessor(preprocessor Preprocessor) *TaskRunner {
	tr.preprocessor = composePreprocessor(preprocessor)
	return tr
}

func (tr *TaskRunner) WithHandler(handler Handler) *TaskRunner {
	tr.responseHandlers["default"] = composeHandler(handler)
	return tr
}

func (tr *TaskRunner) WithCustomHandler(subject string, handler Handler) *TaskRunner {
	tr.responseHandlers[subject] = composeHandler(handler)
	return tr
}

func (tr *TaskRunner) WithPostprocessor(postprocessor Postprocessor) *TaskRunner {
	tr.postprocessor = composePostprocessor(postprocessor)
	return tr
}

func (tr *TaskRunner) WithFinalizer(finalizer common.Finalizer) *TaskRunner {
	tr.finalizer = composeFinalizer(finalizer)
	return tr
}

func (tr *TaskRunner) Run() {
	if tr.responseHandlers["default"] == nil {
		panic("No default handler defined")
	}
	if tr.initializer == nil {
		tr.initializer = composeInitializer(nil)
	}
	if tr.finalizer == nil {
		tr.finalizer = composeFinalizer(nil)
	}

	tr.initializer(tr.ctx)

	go tr.startSubscriber()

	wg.Add(1)
	wg.Wait()

	tr.finalizer(tr.ctx)
}
