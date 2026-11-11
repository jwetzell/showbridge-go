package processor

import (
	"context"
	"errors"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/types"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:         "route.write",
		Title:        "Route Write",
		Description:  "Write to the current route writer if there is one",
		ParamsSchema: nil,
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			return &RouteWrite{config: config}, nil
		},
	})
}

type RouteWrite struct {
	config config.ProcessorConfig
}

func (rw *RouteWrite) Process(next routing.ProcessFunc) routing.ProcessFunc {
	return func(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {

		bytes, ok := types.GetAnyAsByteSlice(wrappedPayload.Payload)
		if !ok {
			return wrappedPayload, errors.New("route.write could not convert payload to byte slice")
		}

		contextWriter, ok := ctx.Value(routing.WriterContextKey).(*routing.ContextWriter)

		if !ok {
			return wrappedPayload, errors.New("route.write input must come from a module that provides a writer")
		}

		if contextWriter.Done() {
			return wrappedPayload, errors.New("route.write response writer has already been written to")
		}

		_, err := contextWriter.Write(bytes)
		if err != nil {
			return wrappedPayload, err
		}

		return next(ctx, wrappedPayload, moduleManager, router)
	}
}
