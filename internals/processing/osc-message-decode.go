package processing

import (
	"context"
	"fmt"

	osc "github.com/jwetzell/osc-go"
)

type OSCMessageDecode struct {
	config ProcessorConfig
}

func (o *OSCMessageDecode) Process(ctx context.Context, payload any) (any, error) {
	payloadBytes, ok := payload.([]byte)

	if !ok {
		return nil, fmt.Errorf("osc.message.decode processor only accepts a []byte payload")
	}

	if len(payloadBytes) == 0 {
		return nil, fmt.Errorf("osc.message.decode processor can't work on empty []byte")
	}

	if payloadBytes[0] != '/' {
		return nil, fmt.Errorf("osc.message.decode processor needs an OSC looking []byte")
	}

	message, err := osc.MessageFromBytes(payloadBytes)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (o *OSCMessageDecode) Type() string {
	return o.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "osc.message.decode",
		New: func(config ProcessorConfig) (Processor, error) {
			return &OSCMessageDecode{config: config}, nil
		},
	})
}
