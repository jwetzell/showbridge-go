package processor

import (
	"context"
	"errors"

	"github.com/jwetzell/showbridge-go/common"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/payload"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "string.decode",
		Title: "Decode String",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &StringDecode{config: config}, nil
		},
	})
}

type StringDecode struct {
	config config.ProcessorConfig
}

func (sd *StringDecode) Process(ctx context.Context, wrappedPayload payload.WrappedPayload, moduleManager common.ModuleManager, inputHandler common.InputHandler) (payload.WrappedPayload, error) {
	payload := wrappedPayload.Payload

	// payload is already a string, return it as is
	_, ok := payload.(string)
	if ok {
		return wrappedPayload, nil
	}

	// payload is convertible to a byte slice
	payloadBytes, ok := common.GetAnyAsByteSlice(payload)
	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("string.decode processor only accepts a []byte")
	}

	payloadMessage := string(payloadBytes)

	wrappedPayload.Payload = payloadMessage
	return wrappedPayload, nil
}
