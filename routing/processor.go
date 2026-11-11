package routing

import (
	"context"
)

type Processor interface {
	Process(next ProcessFunc) ProcessFunc
}

type ProcessFunc func(context.Context, WrappedPayload, ModuleManager, Router) (WrappedPayload, error)
