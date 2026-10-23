package processor

import (
	"context"
	"errors"
	"fmt"

	osc "github.com/jwetzell/osc-go"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/common"
)

type OSCMessageEncode struct {
	config config.ProcessorConfig
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "osc.message.encode",
		Title: "Encode OSC Message",
		New: func(config config.ProcessorConfig) (common.Processor, error) {
			return &OSCMessageEncode{config: config}, nil
		},
	})
}

func (ome *OSCMessageEncode) Process(ctx context.Context, wrappedPayload common.WrappedPayload, moduleManager common.ModuleManager, inputHandler common.InputHandler) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadMessage, ok := common.GetAnyAs[*osc.Message](payload)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("osc.message.encode processor only accepts an *OSCMessage")
	}

	bytes, err := payloadMessage.ToBytes()
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("osc.message.encode processor failed to encode OSCMessage: %w", err)
	}
	wrappedPayload.Payload = bytes
	return wrappedPayload, nil
}
