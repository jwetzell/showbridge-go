package processing

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
)

type NATSMessageEncode struct {
	config ProcessorConfig
}

func (nme *NATSMessageEncode) Process(ctx context.Context, payload any) (any, error) {
	payloadMessage, ok := payload.(*nats.Msg)

	if !ok {
		return nil, fmt.Errorf("nats.message.encode processor only accepts an nats.Msg")
	}

	return payloadMessage.Data, nil
}

func (nme *NATSMessageEncode) Type() string {
	return nme.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "nats.message.encode",
		New: func(config ProcessorConfig) (Processor, error) {
			return &NATSMessageEncode{config: config}, nil
		},
	})
}
