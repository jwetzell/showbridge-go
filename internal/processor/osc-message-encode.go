package processor

import (
	"context"
	"errors"

	osc "github.com/jwetzell/osc-go"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type OSCMessageEncode struct {
	config config.ProcessorConfig
}

func (ome *OSCMessageEncode) Process(ctx context.Context, payload any) (any, error) {
	payloadMessage, ok := payload.(osc.OSCMessage)

	if !ok {
		return nil, errors.New("osc.message.encode processor only accepts an OSCMessage")
	}

	bytes := payloadMessage.ToBytes()
	return bytes, nil
}

func (ome *OSCMessageEncode) Type() string {
	return ome.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "osc.message.encode",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &OSCMessageEncode{config: config}, nil
		},
	})
}
