package routing

import (
	"context"
)

type InputHandler func(ctx context.Context, sourceId string, payload any) (bool, []RouteIOError)

type RouteIOError struct {
	Index        int   `json:"index"`
	ProcessError error `json:"processError"`
}

type Route interface {
	ProcessPayload(ctx context.Context, wrappedPayload WrappedPayload, moduleManager ModuleManager, inputHandler InputHandler) error
	Matches(source string) bool
}
