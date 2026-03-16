package route

import (
	"context"
	"fmt"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type RouteError struct {
	Index  int                `json:"index"`
	Config config.RouteConfig `json:"config"`
	Error  string             `json:"error"`
}
type Route struct {
	input      string
	processors []processor.Processor
}

func NewRoute(config config.RouteConfig) (*Route, error) {
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

	return &Route{input: config.Input, processors: processors}, nil
}

func (r *Route) Input() string {
	return r.input
}

func (r *Route) ProcessPayload(ctx context.Context, payload any) (any, error) {
	wrappedPayload := common.GetWrappedPayload(ctx, payload)
	tracer := otel.Tracer("route")
	processCtx, processSpan := tracer.Start(ctx, "ProcessPayload", trace.WithAttributes(attribute.String("payload.type", fmt.Sprintf("%T", payload))))
	defer processSpan.End()
	for processorIndex, processor := range r.processors {
		processorCtx, processorSpan := otel.Tracer("processor").Start(processCtx, "process", trace.WithAttributes(attribute.Int("processor.index", processorIndex), attribute.String("processor.type", processor.Type())))
		processedPayload, err := processor.Process(processorCtx, wrappedPayload)
		if err != nil {
			processorSpan.SetStatus(codes.Error, "route processor error")
			processorSpan.RecordError(err)
			processorSpan.End()
			processSpan.SetStatus(codes.Error, "route processing error")
			processSpan.RecordError(err)
			return nil, err
		}
		processorSpan.SetStatus(codes.Ok, "processor successful")
		//NOTE(jwetzell) payload has been marked as an end
		if processedPayload.End {
			processSpan.SetStatus(codes.Ok, "route processing terminated early due to processor signal")
			processorSpan.End()
			return processedPayload.Payload, nil
		}
		wrappedPayload = processedPayload
		processorSpan.End()
	}
	processSpan.SetStatus(codes.Ok, "route processing successful")

	return wrappedPayload.Payload, nil
}
