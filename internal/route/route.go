package route

import (
	"context"
	"fmt"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

type routeContextKey string

var RouterContextKey routeContextKey = routeContextKey("router")
var SourceContextKey routeContextKey = routeContextKey("source")

type RouteError struct {
	Index  int
	Config config.RouteConfig
	Error  error
}

type RouteIOError struct {
	Index int
	Error error
}

type RouteIO interface {
	HandleInput(sourceId string, payload any) []RouteIOError
	HandleOutput(ctx context.Context, destinationId string, payload any) error
}

type Route interface {
	Input() string
	Output() string
	ProcessPayload(ctx context.Context, payload any) (any, error)
}

type ProcessorRoute struct {
	input      string
	processors []processor.Processor
	output     string
}

func NewRoute(config config.RouteConfig) (Route, error) {
	processors := []processor.Processor{}

	if len(config.Processors) > 0 {
		for _, processorDecl := range config.Processors {
			processorInfo, ok := processor.ProcessorRegistry[processorDecl.Type]
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

	return &ProcessorRoute{input: config.Input, processors: processors, output: config.Output}, nil
}

func (r *ProcessorRoute) Input() string {
	return r.input
}

func (r *ProcessorRoute) Output() string {
	return r.output
}

func (r *ProcessorRoute) ProcessPayload(ctx context.Context, payload any) (any, error) {
	for _, processor := range r.processors {
		processedPayload, err := processor.Process(ctx, payload)
		if err != nil {
			return nil, err
		}
		//NOTE(jwetzell) nil payload will result in the route being "terminated"
		if processedPayload == nil {
			return nil, nil
		}
		payload = processedPayload
	}

	return payload, nil
}
