package route

import (
	"context"
	"fmt"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

type Route struct {
	input      string
	processors []processor.Processor
}

func NewRoute(config config.RouteConfig) (*Route, error) {
	processors := []processor.Processor{}

	if len(config.Processors) > 0 {
		for _, processorDecl := range config.Processors {
			processorInfo, ok := processor.GetProcessorRegistration(processorDecl.Type)
			if !ok {
				return nil, fmt.Errorf("problem loading processor registration for processor type: %s", processorDecl.Type)
			}

			processor, err := processorInfo.New(processorDecl)
			if err != nil {
				return nil, err
			}
			processors = append(processors, processor)
		}
	}

	return &Route{input: config.Input, processors: processors}, nil
}

func (r *Route) Input() string {
	return r.input
}

func (r *Route) ProcessPayload(ctx context.Context, wrappedPayload common.WrappedPayload) (any, error) {
	for _, processor := range r.processors {
		processedPayload, err := processor.Process(ctx, wrappedPayload)
		if err != nil {
			return nil, err
		}
		//NOTE(jwetzell) payload has been marked as an end
		if processedPayload.End {
			return processedPayload.Payload, nil
		}
		wrappedPayload = processedPayload
	}

	return wrappedPayload.Payload, nil
}
