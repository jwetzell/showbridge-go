package test

import (
	"context"

	"github.com/jwetzell/showbridge-go/internal/common"
)

type TestRouter struct {
}

func (r *TestRouter) HandleInput(ctx context.Context, sourceId string, payload any) (bool, []common.RouteIOError) {
	return false, nil
}

func (r *TestRouter) HandleOutput(ctx context.Context, destinationId string, payload any) error {
	return nil
}

func GetNewTestRouter() *TestRouter {
	return &TestRouter{}
}
