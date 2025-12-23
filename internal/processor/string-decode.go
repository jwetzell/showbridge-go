package processor

import (
	"context"
	"errors"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type StringDecode struct {
	config config.ProcessorConfig
}

func (sd *StringDecode) Process(ctx context.Context, payload any) (any, error) {
	payloadBytes, ok := payload.([]byte)

	if !ok {
		return nil, errors.New("string.decode processor only accepts a []byte")
	}

	payloadMessage := string(payloadBytes)

	return payloadMessage, nil
}

func (sd *StringDecode) Type() string {
	return sd.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "string.decode",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &StringDecode{config: config}, nil
		},
	})
}
