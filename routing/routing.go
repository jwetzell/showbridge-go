package routing

import (
	"context"
)

type RouteIOError struct {
	Index        int   `json:"index"`
	ProcessError error `json:"processError"`
}

type Route interface {
	ProcessPayload(ctx context.Context, wrappedPayload WrappedPayload, moduleManager ModuleManager, router Router) error
	Matches(source string) bool
}

type Router interface {
	HandleInput(ctx context.Context, sourceId string, inputPayload any) (bool, []RouteIOError)
}
