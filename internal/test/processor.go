package test

import (
	"context"

	"github.com/jwetzell/showbridge-go/common"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/payload"
)

type TestProcessor struct {
	Config config.ProcessorConfig
}

func (p *TestProcessor) Process(ctx context.Context, wrappedPayload payload.WrappedPayload, moduleManager common.ModuleManager, inputHandler common.InputHandler) (payload.WrappedPayload, error) {
	return wrappedPayload, nil
}
