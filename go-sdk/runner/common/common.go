package common

import (
	"github.com/konstellation-io/kre-runners/go-sdk/v1/context"
	"google.golang.org/protobuf/types/known/anypb"
)

type Task func(ctx context.KaiContext)

type Initializer Task

type Finalizer Task

type Handler func(ctx context.KaiContext, response *anypb.Any) error
