package routing

import (
	"context"
)

type Processor interface {
	Process(context.Context, WrappedPayload, ModuleManager, Router) (WrappedPayload, error)
}
