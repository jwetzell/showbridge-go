package processor

import (
	"context"
	"fmt"

	osc "github.com/jwetzell/osc-go"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type OSCMessageEncode struct {
	config config.ProcessorConfig
}

func (o *OSCMessageEncode) Process(ctx context.Context, payload any) (any, error) {
	payloadMessage, ok := payload.(osc.OSCMessage)

	if !ok {
		return nil, fmt.Errorf("osc.message.encode processor only accepts an OSCMessage")
	}

	bytes := payloadMessage.ToBytes()
	return bytes, nil
}

func (o *OSCMessageEncode) Type() string {
	return o.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "osc.message.encode",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &OSCMessageEncode{config: config}, nil
		},
	})
}
