package processor

import (
	"context"
	"errors"

	"github.com/jwetzell/showbridge-go/common"
	"github.com/jwetzell/showbridge-go/config"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "string.encode",
		Title: "Encode String",
		New: func(config config.ProcessorConfig) (common.Processor, error) {
			return &StringEncode{config: config}, nil
		},
	})
}

type StringEncode struct {
	config config.ProcessorConfig
}

func (se *StringEncode) Process(ctx context.Context, wrappedPayload common.WrappedPayload, moduleManager common.ModuleManager, inputHandler common.InputHandler) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadString, ok := common.GetAnyAs[string](payload)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("string.encode processor only accepts a string")
	}

	wrappedPayload.Payload = []byte(payloadString)

	return wrappedPayload, nil
}
