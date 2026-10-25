package showbridge

import (
	"context"
	"fmt"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/processor"
	"github.com/jwetzell/showbridge-go/routing"
)

type Route struct {
	id         string
	input      string
	processors []routing.Processor
}

func (r *Route) Matches(source string) bool {
	// TODO(jwetzell): maybe support more complex matching logic in the future? regex? wildcards?
	return r.input == source
}

func (r *Route) ProcessPayload(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) error {
	for processorIndex, processor := range r.processors {
		processedPayload, err := processor.Process(ctx, wrappedPayload, moduleManager, router)
		if err != nil {
			return fmt.Errorf("processor[%d] error: %w", processorIndex, err)
		}
		// NOTE(jwetzell) payload has been marked as an end without error
		if processedPayload.End {
			return nil
		}
		wrappedPayload = processedPayload
	}

	return nil
}

func NewRoute(routeDecl config.RouteConfig) (*Route, error) {
	processors := []routing.Processor{}

	if len(routeDecl.Processors) > 0 {
		for _, processorDecl := range routeDecl.Processors {
			processorInfo, ok := processor.GetProcessorRegistration(processorDecl.Type)
			if !ok {
				return nil, fmt.Errorf("problem loading processor registration for processor type: %s", processorDecl.Type)
			}

			proc, err := processorInfo.New(processorDecl)
			if err != nil {
				return nil, err
			}
			processors = append(processors, proc)
		}
	}
	routeInstance := &Route{id: routeDecl.Id, input: routeDecl.Input, processors: processors}
	return routeInstance, nil
}
