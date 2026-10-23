package common

import (
	"context"
)

type Processor interface {
	Process(context.Context, WrappedPayload, ModuleManager, InputHandler) (WrappedPayload, error)
}
