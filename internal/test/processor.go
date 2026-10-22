package test

import (
	"context"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/common"
)

type TestProcessor struct {
	Config config.ProcessorConfig
}

func (p *TestProcessor) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	return wrappedPayload, nil
}
