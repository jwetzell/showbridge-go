package route

import (
	"context"
	"fmt"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

type route struct {
	id         string
	input      string
	processors []common.Processor
}

func NewRoute(config config.RouteConfig) (common.Route, error) {
	processors := []common.Processor{}

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

	return &route{id: config.Id, input: config.Input, processors: processors}, nil
}

func (r *route) Matches(source string) bool {
	// TODO(jwetzell): maybe support more complex matching logic in the future? regex? wildcards?
	return r.input == source
}

func (r *route) ProcessPayload(ctx context.Context, wrappedPayload common.WrappedPayload) (any, error) {
	for processorIndex, processor := range r.processors {
		processedPayload, err := processor.Process(ctx, wrappedPayload)
		if err != nil {
			return nil, fmt.Errorf("processor[%d] error: %w", processorIndex, err)
		}
		// NOTE(jwetzell) payload has been marked as an end without error
		if processedPayload.End {
			return processedPayload.Payload, nil
		}
		wrappedPayload = processedPayload
	}

	return wrappedPayload.Payload, nil
}
