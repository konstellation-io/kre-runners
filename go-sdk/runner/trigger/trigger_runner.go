package trigger

import (
	"github.com/go-logr/logr"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/context"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/runner/common"
	"github.com/nats-io/nats.go"
	"sync"
)

//type Runner common.Task

type Runner func(ctx context.KaiContext, handlers map[string]ResponseHandler)

type ResponseHandler common.Handler

type TriggerRunner struct {
	ctx              context.KaiContext
	nats             *nats.Conn
	jetstream        nats.JetStreamContext
	responseHandler  ResponseHandler
	responseHandlers map[string]ResponseHandler
	initializer      common.Initializer
	runner           Runner
	finalizer        common.Finalizer
}

var wg sync.WaitGroup

func NewTriggerRunner(logger logr.Logger, nats *nats.Conn, jetstream nats.JetStreamContext) *TriggerRunner {
	return &TriggerRunner{
		ctx:              context.NewContext(logger.WithName("[Trigger]"), nats, jetstream),
		nats:             nats,
		jetstream:        jetstream,
		responseHandlers: make(map[string]ResponseHandler),
	}
}

func (tr *TriggerRunner) WithInitializer(initializer common.Initializer) *TriggerRunner {
	tr.initializer = composeInitializer(initializer)
	return tr
}

func (tr *TriggerRunner) WithRunner(runner Runner) *TriggerRunner {
	tr.runner = composeRunner(runner)
	return tr
}

func (tr *TriggerRunner) WithFinalizer(finalizer common.Finalizer) *TriggerRunner {
	tr.finalizer = composeFinalizer(finalizer)
	return tr
}

func (tr *TriggerRunner) Run() {
	// Check required fields are initialized
	if tr.runner == nil {
		panic("Runner function is required")
	}
	if tr.initializer == nil {
		tr.initializer = composeInitializer(nil)
	}

	tr.responseHandler = composeResponseHandler(tr.responseHandlers)
	if tr.finalizer == nil {
		tr.finalizer = composeFinalizer(nil)
	}

	tr.initializer(tr.ctx)

	go tr.runner(tr.ctx, tr.responseHandlers)
	// SendOutput() --> event
	// exit node responds to event
	// goroutine to wait for the response (passing the channel)

	go tr.startSubscriber()
	// Calls the handler
	// Respond to the channel with the execution ID

	wg.Add(2)
	wg.Wait()

	tr.finalizer(tr.ctx)
}
