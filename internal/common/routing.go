package common

import "context"

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
