package processor

import (
	"context"
	"errors"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type MQTTMessageEncode struct {
	config config.ProcessorConfig
}

func (mme *MQTTMessageEncode) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadMessage, ok := common.GetAnyAs[mqtt.Message](payload)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("mqtt.message.encode processor only accepts an mqtt.Message")
	}
	wrappedPayload.Payload = payloadMessage.Payload()
	return wrappedPayload, nil
}

func (mme *MQTTMessageEncode) Type() string {
	return mme.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "mqtt.message.encode",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &MQTTMessageEncode{config: config}, nil
		},
	})
}
