package processor

import (
	"context"
	"reflect"

	"github.com/jwetzell/showbridge-go/common"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/payload"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "filter.change",
		Title: "Filter On Change",
		New: func(config config.ProcessorConfig) (common.Processor, error) {
			return &FilterChange{config: config}, nil
		},
	})
}

type FilterChange struct {
	config   config.ProcessorConfig
	previous any
}

func (fc *FilterChange) Process(ctx context.Context, wrappedPayload payload.WrappedPayload, moduleManager common.ModuleManager, inputHandler common.InputHandler) (payload.WrappedPayload, error) {
	payload := wrappedPayload.Payload

	if reflect.DeepEqual(payload, fc.previous) {
		wrappedPayload.End = true
		return wrappedPayload, nil
	}
	fc.previous = payload

	return wrappedPayload, nil
}
