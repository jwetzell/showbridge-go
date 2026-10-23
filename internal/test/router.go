package test

import (
	"context"

	"github.com/jwetzell/showbridge-go/common"
)

type TestRouter struct {
}

func (r *TestRouter) HandleInput(ctx context.Context, sourceId string, payload any) (bool, []common.RouteIOError) {
	return false, nil
}

func GetNewTestRouter() *TestRouter {
	return &TestRouter{}
}
