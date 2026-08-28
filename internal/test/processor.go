package test

import (
	"context"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type TestProcessor struct {
	Config config.ProcessorConfig
}

func (p *TestProcessor) Id() string {
	return p.Config.Id
}

func (p *TestProcessor) Type() string {
	return p.Config.Type
}
func (p *TestProcessor) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	return wrappedPayload, nil
}
