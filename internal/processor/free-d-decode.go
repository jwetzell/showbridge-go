package processor

import (
	"context"
	"errors"

	freeD "github.com/jwetzell/free-d-go"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type FreeDDecode struct {
	config config.ProcessorConfig
}

func (fdd *FreeDDecode) Process(ctx context.Context, payload any) (any, error) {
	payloadBytes, ok := payload.([]byte)

	if !ok {
		return nil, errors.New("freed.decode processor only accepts a []byte")
	}

	payloadMessage, err := freeD.Decode(payloadBytes)
	if err != nil {
		return nil, err
	}
	return payloadMessage, nil
}

func (fdd *FreeDDecode) Type() string {
	return fdd.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "freed.decode",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &FreeDDecode{config: config}, nil
		},
	})
}
