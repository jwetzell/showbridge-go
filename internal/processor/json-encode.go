package processor

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type JsonEncode struct {
	config config.ProcessorConfig
}

func (je *JsonEncode) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	var payloadBuffer bytes.Buffer

	err := json.NewEncoder(&payloadBuffer).Encode(payload)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	payloadBytes := payloadBuffer.Bytes()

	payloadBytes = payloadBytes[0 : len(payloadBytes)-1]

	wrappedPayload.Payload = payloadBytes
	return wrappedPayload, nil
}

func (je *JsonEncode) Type() string {
	return je.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "json.encode",
		Title: "Encode JSON",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &JsonEncode{config: config}, nil
		},
	})
}
