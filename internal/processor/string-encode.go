package processor

import (
	"context"
	"errors"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type StringEncode struct {
	config config.ProcessorConfig
}

func (se *StringEncode) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadString, ok := common.GetAnyAs[string](payload)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("string.encode processor only accepts a string")
	}

	wrappedPayload.Payload = []byte(payloadString)

	return wrappedPayload, nil
}

func (se *StringEncode) Type() string {
	return se.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "string.encode",
		Title: "Encode String",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &StringEncode{config: config}, nil
		},
	})
}
