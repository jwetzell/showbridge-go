package common

import (
	"context"

	"github.com/jwetzell/showbridge-go/payload"
)

type InputHandler func(ctx context.Context, sourceId string, payload any) (bool, []RouteIOError)

type RouteIOError struct {
	Index        int   `json:"index"`
	ProcessError error `json:"processError"`
}

type Route interface {
	ProcessPayload(ctx context.Context, wrappedPayload payload.WrappedPayload, moduleManager ModuleManager, inputHandler InputHandler) (any, error)
	Matches(source string) bool
}
