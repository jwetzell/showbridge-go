package processing

import (
	"context"
	"fmt"
)

type StringDecode struct {
	config ProcessorConfig
}

func (sd *StringDecode) Process(ctx context.Context, payload any) (any, error) {
	payloadBytes, ok := payload.([]byte)

	if !ok {
		return nil, fmt.Errorf("string.decode processor only accepts a []byte")
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
		New: func(config ProcessorConfig) (Processor, error) {
			return &StringDecode{config: config}, nil
		},
	})
}
