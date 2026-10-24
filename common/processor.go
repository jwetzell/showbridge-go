package common

import (
	"context"

	"github.com/jwetzell/showbridge-go/payload"
)

type Processor interface {
	Process(context.Context, payload.WrappedPayload, ModuleManager, InputHandler) (payload.WrappedPayload, error)
}
