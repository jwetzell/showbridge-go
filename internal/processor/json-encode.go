package processor

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type JsonEncode struct {
	config config.ProcessorConfig
}

func (je *JsonEncode) Process(ctx context.Context, payload any) (any, error) {
	var payloadBuffer bytes.Buffer

	err := json.NewEncoder(&payloadBuffer).Encode(payload)

	if err != nil {
		return nil, err
	}

	payloadBytes := payloadBuffer.Bytes()

	payloadBytes = payloadBytes[0 : len(payloadBytes)-1]

	return payloadBytes, nil
}

func (je *JsonEncode) Type() string {
	return je.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "json.encode",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &JsonEncode{config: config}, nil
		},
	})
}
