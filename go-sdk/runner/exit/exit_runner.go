package exit

import (
	"github.com/go-logr/logr"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/context"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/runner/common"
	"github.com/nats-io/nats.go"
	"sync"
)

type Preprocessor common.Handler

type Handler common.Handler

type Postprocessor common.Handler

type ExitRunner struct {
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

func NewExitRunner(logger logr.Logger, nats *nats.Conn, jetstream nats.JetStreamContext) *ExitRunner {
	return &ExitRunner{
		ctx:              context.NewContext(logger.WithName("[Exit]"), nats, jetstream),
		nats:             nats,
		jetstream:        jetstream,
		responseHandlers: make(map[string]Handler),
	}
}

func (er *ExitRunner) WithInitializer(initializer common.Initializer) *ExitRunner {
	er.initializer = composeInitializer(initializer)
	return er
}

func (er *ExitRunner) WithPreprocessor(preprocessor Preprocessor) *ExitRunner {
	er.preprocessor = composePreprocessor(preprocessor)
	return er
}

func (er *ExitRunner) WithHandler(handler Handler) *ExitRunner {
	er.responseHandlers["default"] = composeHandler(handler)
	return er
}

func (er *ExitRunner) WithPostprocessor(postprocessor Postprocessor) *ExitRunner {
	er.postprocessor = composePostprocessor(postprocessor)
	return er
}

func (er *ExitRunner) WithFinalizer(finalizer common.Finalizer) *ExitRunner {
	er.finalizer = composeFinalizer(finalizer)
	return er
}

func (er *ExitRunner) Run() {
	if er.responseHandlers["default"] == nil {
		panic("No default handler defined")
	}
	if er.initializer == nil {
		er.initializer = composeInitializer(nil)
	}
	if er.finalizer == nil {
		er.finalizer = composeFinalizer(nil)
	}

	er.initializer(er.ctx)

	go er.startSubscriber()

	wg.Add(1)
	wg.Wait()

	er.finalizer(er.ctx)
}
