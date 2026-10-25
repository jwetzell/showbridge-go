package showbridge

import (
	"context"
	"fmt"

	"github.com/jwetzell/showbridge-go/routing"
)

type route struct {
	id         string
	input      string
	processors []routing.Processor
}

func (r *route) Matches(source string) bool {
	// TODO(jwetzell): maybe support more complex matching logic in the future? regex? wildcards?
	return r.input == source
}

func (r *route) ProcessPayload(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, inputHandler routing.InputHandler) (any, error) {
	for processorIndex, processor := range r.processors {
		processedPayload, err := processor.Process(ctx, wrappedPayload, moduleManager, inputHandler)
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
