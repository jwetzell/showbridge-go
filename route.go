package showbridge

import (
	"context"
	"fmt"
	"slices"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/processor"
	"github.com/jwetzell/showbridge-go/routing"
)

type Route struct {
	id          string
	input       string
	processors  []routing.Processor
	processFunc routing.ProcessFunc
}

func (r *Route) Matches(source string) bool {
	// TODO(jwetzell): maybe support more complex matching logic in the future? regex? wildcards?
	return r.input == source
}

func (r *Route) ProcessPayload(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) error {

	if r.processFunc == nil {
		return nil
	}

	_, err := r.processFunc(ctx, wrappedPayload, moduleManager, router)
	if err != nil {
		return fmt.Errorf("processing error: %w", err)
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

	routeInstance.processFunc = func(ctx context.Context, wp routing.WrappedPayload, mm routing.ModuleManager, r routing.Router) (routing.WrappedPayload, error) {
		return wp, nil
	}

	for _, processor := range slices.Backward(processors) {
		routeInstance.processFunc = processor.Process(routeInstance.processFunc)
	}

	return routeInstance, nil
}
