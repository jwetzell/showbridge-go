package common

import "context"

type contextKey string

const RouterContextKey contextKey = contextKey("router")
const SourceContextKey contextKey = contextKey("source")
const ModulesContextKey contextKey = contextKey("modules")
const SenderContextKey contextKey = contextKey("sender")

type RouteIO interface {
	HandleInput(ctx context.Context, sourceId string, payload any) (bool, []RouteIOError)
	HandleOutput(ctx context.Context, destinationId string, payload any) error
}

type RouteIOError struct {
	Index        int
	OutputError  error
	ProcessError error
	InputError   error
}
