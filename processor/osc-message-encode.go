package processor

import (
	"context"
	"errors"
	"fmt"

	osc "github.com/jwetzell/osc-go"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/routing"
)

type OSCMessageEncode struct {
	config config.ProcessorConfig
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "osc.message.encode",
		Title: "Encode OSC Message",
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			return &OSCMessageEncode{config: config}, nil
		},
	})
}

func (ome *OSCMessageEncode) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadMessage, ok := payload.(*osc.Message)

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
