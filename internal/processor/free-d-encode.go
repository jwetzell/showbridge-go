package processor

import (
	"context"
	"fmt"

	freeD "github.com/jwetzell/free-d-go"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type FreeDEncode struct {
	config config.ProcessorConfig
}

func (fde *FreeDEncode) Process(ctx context.Context, payload any) (any, error) {
	payloadPosition, ok := payload.(freeD.FreeDPosition)

	if !ok {
		return nil, fmt.Errorf("freed.decode processor only accepts a FreeDEncode")
	}

	payloadBytes := freeD.Encode(payloadPosition)
	return payloadBytes, nil
}

func (fde *FreeDEncode) Type() string {
	return fde.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "freed.encode",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &FreeDEncode{config: config}, nil
		},
	})
}
