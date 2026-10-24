package processor

import (
	"context"
	"errors"

	freeD "github.com/jwetzell/free-d-go"
	"github.com/jwetzell/showbridge-go/common"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/payload"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "freed.encode",
		Title: "Encode FreeD",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &FreeDEncode{config: config}, nil
		},
	})
}

type FreeDEncode struct {
	config config.ProcessorConfig
}

func (fe *FreeDEncode) Process(ctx context.Context, wrappedPayload payload.WrappedPayload, moduleManager common.ModuleManager, inputHandler common.InputHandler) (payload.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadPosition, ok := payload.(freeD.FreeDPosition)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("freed.encode processor only accepts a FreeDPosition")
	}

	payloadBytes := freeD.Encode(payloadPosition)

	wrappedPayload.Payload = payloadBytes
	return wrappedPayload, nil
}
