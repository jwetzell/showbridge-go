package common

import (
	"context"
)

type RouteIO interface {
	HandleInput(ctx context.Context, sourceId string, payload any) (bool, []RouteIOError)
}

type RouteIOError struct {
	Index        int   `json:"index"`
	ProcessError error `json:"processError"`
}
