package processor

import (
	"context"
	"reflect"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type FilterUnique struct {
	config   config.ProcessorConfig
	previous any
}

func (fr *FilterUnique) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload

	if reflect.DeepEqual(payload, fr.previous) {
		wrappedPayload.End = true
		return wrappedPayload, nil
	}
	fr.previous = payload

	return wrappedPayload, nil
}

func (fr *FilterUnique) Type() string {
	return fr.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "filter.unique",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &FilterUnique{config: config}, nil
		},
	})
}
