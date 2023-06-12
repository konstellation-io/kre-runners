package task

import (
	"github.com/konstellation-io/kre-runners/go-sdk/v1/context"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/runner/common"
	"google.golang.org/protobuf/types/known/anypb"
)

func composeInitializer(initializer common.Initializer) common.Initializer {
	return func(ctx context.KaiContext) {
		ctx.Logger.V(1).Info("Initializing TaskRunner...")
		if initializer != nil {
			initializer(ctx)
		}
	}
}

func composePreprocessor(preprocessor Preprocessor) Preprocessor {
	return func(ctx context.KaiContext, response *anypb.Any) error {
		ctx.Logger.V(1).Info("Preprocessing TaskRunner...")

		if preprocessor != nil {
			return preprocessor(ctx, response)
		}

		return nil
	}
}

func composeHandler(handler Handler) Handler {
	return func(ctx context.KaiContext, response *anypb.Any) error {
		ctx.Logger.V(1).Info("Handling TaskRunner...")

		if handler != nil {
			return handler(ctx, response)
		}

		return nil
	}
}

func composePostprocessor(postprocessor Postprocessor) Postprocessor {
	return func(ctx context.KaiContext, response *anypb.Any) error {
		ctx.Logger.V(1).Info("Postprocessing TaskRunner...")

		if postprocessor != nil {
			return postprocessor(ctx, response)
		}

		return nil
	}
}

func composeFinalizer(finalizer common.Finalizer) common.Finalizer {
	return func(ctx context.KaiContext) {
		ctx.Logger.V(1).Info("Finalizing TaskRunner...")

		if finalizer != nil {
			finalizer(ctx)
		}
	}
}
