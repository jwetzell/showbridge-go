package processor

import (
	"context"
	"errors"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:        "string.encode",
		Title:       "Encode String",
		Description: "Encodes a string as bytes",
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			return &StringEncode{config: config}, nil
		},
	})
}

type StringEncode struct {
	config config.ProcessorConfig
}

func (se *StringEncode) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadString, ok := payload.(string)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("string.encode processor only accepts a string")
	}

	wrappedPayload.Payload = []byte(payloadString)

	return wrappedPayload, nil
}
