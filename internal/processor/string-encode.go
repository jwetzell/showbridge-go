package processor

import (
	"context"
	"errors"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type StringEncode struct {
	config config.ProcessorConfig
}

func (se *StringEncode) Process(ctx context.Context, payload any) (any, error) {
	payloadString, ok := payload.(string)

	if !ok {
		return nil, errors.New("string.encode processor only accepts a string")
	}

	payloadBytes := []byte(payloadString)

	return payloadBytes, nil
}

func (se *StringEncode) Type() string {
	return se.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "string.encode",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &StringEncode{config: config}, nil
		},
	})
}
