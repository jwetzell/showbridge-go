package test

import (
	"context"

	"github.com/jwetzell/showbridge-go/internal/common"
)

type TestProcessor struct {
}

func (p *TestProcessor) Type() string {
	return "test"
}
func (p *TestProcessor) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	return wrappedPayload, nil
}
