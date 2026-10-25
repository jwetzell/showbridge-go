package processor

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "json.encode",
		Title: "Encode JSON",
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			return &JsonEncode{config: config}, nil
		},
	})
}

type JsonEncode struct {
	config config.ProcessorConfig
}

func (je *JsonEncode) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, inputHandler routing.InputHandler) (routing.WrappedPayload, error) {
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
