package processor

import (
	"context"
	"reflect"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type FilterChange struct {
	config   config.ProcessorConfig
	previous any
}

func (fc *FilterChange) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload

	if reflect.DeepEqual(payload, fc.previous) {
		wrappedPayload.End = true
		return wrappedPayload, nil
	}
	fc.previous = payload

	return wrappedPayload, nil
}

func (fc *FilterChange) Type() string {
	return fc.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "filter.change",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &FilterChange{config: config}, nil
		},
	})
}
