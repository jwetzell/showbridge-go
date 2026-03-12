package common

import "context"

type RouteIO interface {
	HandleInput(ctx context.Context, sourceId string, payload any) (bool, []RouteIOError)
	HandleOutput(ctx context.Context, destinationId string, payload any) error
}

type RouteIOError struct {
	Index        int   `json:"index"`
	OutputError  error `json:"outputError"`
	ProcessError error `json:"processError"`
	InputError   error `json:"inputError"`
}
