package route

import (
	"context"
	"fmt"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
	Index        int
	OutputError  error
	ProcessError error
	InputError   error
}

type RouteIO interface {
	HandleInput(ctx context.Context, sourceId string, payload any) (bool, []RouteIOError)
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
	parentSpan := trace.SpanFromContext(ctx)
	tracer := parentSpan.TracerProvider().Tracer("route.ProcessPayload")
	processCtx, processSpan := tracer.Start(ctx, "route.process")
	defer processSpan.End()
	for processorIndex, processor := range r.processors {
		processorCtx, processorSpan := tracer.Start(processCtx, "route.processor", trace.WithAttributes(attribute.Int("processor.index", processorIndex), attribute.String("processor.type", processor.Type())))
		processedPayload, err := processor.Process(processorCtx, payload)
		if err != nil {
			processorSpan.SetStatus(codes.Error, "route processor error")
			processorSpan.RecordError(err)
			processorSpan.End()
			processSpan.SetStatus(codes.Error, "route processing error")
			processSpan.RecordError(err)
			return nil, err
		}
		processorSpan.SetStatus(codes.Ok, "processor successful")
		//NOTE(jwetzell) nil payload will result in the route being "terminated"
		if processedPayload == nil {
			processSpan.SetStatus(codes.Ok, "route processing terminated early due to nil payload")
			processorSpan.End()
			return nil, nil
		}
		payload = processedPayload
		processorSpan.End()
	}
	processSpan.SetStatus(codes.Ok, "route processing successful")

	return payload, nil
}
