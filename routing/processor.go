package routing

import (
	"context"
)

type Processor interface {
	Process(next ProcessFunc) ProcessFunc
}

type ProcessFunc func(context.Context, WrappedPayload, ModuleManager, Router) (WrappedPayload, error)

func NewLinearProcessFunc(next ProcessFunc, process ProcessFunc) ProcessFunc {
	return func(ctx context.Context, wrappedPayload WrappedPayload, moduleManager ModuleManager, router Router) (WrappedPayload, error) {
		wp, err := process(ctx, wrappedPayload, moduleManager, router)

		if err != nil {
			return wrappedPayload, err
		}
		if wp.End {
			return wp, nil
		}
		return next(ctx, wp, moduleManager, router)
	}
}
