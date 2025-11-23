package processing

import (
	"context"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTMessageEncode struct {
	config ProcessorConfig
}

func (mme *MQTTMessageEncode) Process(ctx context.Context, payload any) (any, error) {
	payloadMessage, ok := payload.(mqtt.Message)

	if !ok {
		return nil, fmt.Errorf("mqtt.message.encode processor only accepts an mqtt.Message")
	}

	return payloadMessage.Payload(), nil
}

func (mme *MQTTMessageEncode) Type() string {
	return mme.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "mqtt.message.encode",
		New: func(config ProcessorConfig) (Processor, error) {
			return &MQTTMessageEncode{config: config}, nil
		},
	})
}
