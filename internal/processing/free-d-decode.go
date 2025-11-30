package processing

import (
	"context"
	"fmt"
	"log/slog"

	freeD "github.com/jwetzell/free-d-go"
)

type FreeDDecode struct {
	config ProcessorConfig
}

func (fdd *FreeDDecode) Process(ctx context.Context, payload any) (any, error) {
	payloadBytes, ok := payload.([]byte)

	if !ok {
		return nil, fmt.Errorf("freed.decode processor only accepts a []byte")
	}

	payloadMessage, err := freeD.Decode(payloadBytes)
	if err != nil {
		slog.Error("error decoding", "err", err)
	}
	return payloadMessage, nil
}

func (fdd *FreeDDecode) Type() string {
	return fdd.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "freed.decode",
		New: func(config ProcessorConfig) (Processor, error) {
			return &FreeDDecode{config: config}, nil
		},
	})
}
