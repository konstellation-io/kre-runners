package trigger

import (
	"github.com/konstellation-io/kre-runners/go-sdk/v1/context"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/runner/common"
	"google.golang.org/protobuf/types/known/anypb"
	"os"
	"os/signal"
	"syscall"
)

func composeInitializer(userInitializer common.Initializer) common.Initializer {
	return func(ctx context.KaiContext) {
		ctx.Logger.V(1).Info("Initializing TriggerRunner...")

		if userInitializer != nil {
			userInitializer(ctx)
		}
	}
}

func composeRunner(userRunner Runner) Runner {
	return func(ctx context.KaiContext, handlers map[string]ResponseHandler) {
		defer wg.Done()

		ctx.Logger.V(1).Info("Running TriggerRunner...")

		if userRunner != nil {
			userRunner(ctx, handlers)
		}

		// Handle sigterm and await termChan signal
		termChan := make(chan os.Signal, 1)
		signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
		<-termChan

		// Handle shutdown
		ctx.Logger.Info("Shutdown signal received")
	}
}

func composeResponseHandler(handlers map[string]ResponseHandler) ResponseHandler {
	return func(ctx context.KaiContext, response *anypb.Any) error {
		ctx.Logger.V(1).Info("Running response handler...")

		// Unmarshal response to a KaiNatsMessage type
		ctx.Logger.Info("Message received", "RequestID", ctx.GetRequestID())

		responseHandler := handlers[ctx.GetRequestID()]

		if responseHandler != nil {
			return responseHandler(ctx, response)
		}

		ctx.Logger.V(1).Info("No handler found for the message",
			"RequestID", ctx.GetRequestID())

		return nil
	}
}

func composeFinalizer(userFinalizer common.Finalizer) common.Finalizer {
	return func(ctx context.KaiContext) {
		ctx.Logger.V(1).Info("Finalizing TriggerRunner...")

		if userFinalizer != nil {
			userFinalizer(ctx)
		}
	}
}
