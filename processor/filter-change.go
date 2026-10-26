package processor

import (
	"context"
	"reflect"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:        "filter.change",
		Title:       "Filter On Change",
		Description: "Filters out payloads that have not changed since the last one",
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			return &FilterChange{config: config}, nil
		},
	})
}

type FilterChange struct {
	config   config.ProcessorConfig
	previous any
}

func (fc *FilterChange) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
	payload := wrappedPayload.Payload

	if reflect.DeepEqual(payload, fc.previous) {
		wrappedPayload.End = true
		return wrappedPayload, nil
	}
	fc.previous = payload

	return wrappedPayload, nil
}
