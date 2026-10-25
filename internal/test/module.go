package test

import (
	"context"

	"github.com/jwetzell/showbridge-go/common"
	_ "modernc.org/sqlite"
)

func NewTestModule(id string) *TestModule {
	return &TestModule{
		id: id,
	}
}

type TestModule struct {
	id string
}

func (m *TestModule) Start(ctx context.Context, inputHandler common.InputHandler) error {
	<-ctx.Done()
	return nil
}

func (m *TestModule) Stop() {}
