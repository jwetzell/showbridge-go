package processor

import (
	"context"
	"errors"

	freeD "github.com/jwetzell/free-d-go"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "freed.encode",
		Title: "Encode FreeD",
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			return &FreeDEncode{config: config}, nil
		},
	})
}

type FreeDEncode struct {
	config config.ProcessorConfig
}

func (fe *FreeDEncode) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, inputHandler routing.InputHandler) (routing.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadPosition, ok := payload.(freeD.Position)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("freed.encode processor only accepts a FreeD Position")
	}

	payloadBytes := payloadPosition.Encode()

	wrappedPayload.Payload = payloadBytes
	return wrappedPayload, nil
}
