package processor

import (
	"context"
	"fmt"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type MQTTMessage struct {
	Topic    string
	QoS      byte
	Payload  any
	Retained bool
}

type MQTTMessageCreate struct {
	config   config.ProcessorConfig
	Topic    string
	QoS      byte
	Retained bool
	Payload  any
}

func (mmc *MQTTMessageCreate) Process(ctx context.Context, payload any) (any, error) {

	message := MQTTMessage{
		Topic:    mmc.Topic,
		QoS:      mmc.QoS,
		Retained: mmc.Retained,
		Payload:  mmc.Payload,
	}

	return message, nil
}

func (mmc *MQTTMessageCreate) Type() string {
	return mmc.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "mqtt.message.create",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params
			topic, ok := params["topic"]

			if !ok {
				return nil, fmt.Errorf("mqtt.message.create requires an topic parameter")
			}

			topicString, ok := topic.(string)

			if !ok {
				return nil, fmt.Errorf("mqtt.message.create topic must be a string")
			}

			qos, ok := params["qos"]

			if !ok {
				return nil, fmt.Errorf("mqtt.message.create requires an qos parameter")
			}

			qosByte, ok := qos.(float64)

			if !ok {
				return nil, fmt.Errorf("mqtt.message.create qos must be a number")
			}

			retained, ok := params["retained"]

			if !ok {
				return nil, fmt.Errorf("mqtt.message.create requires an retained parameter")
			}

			retainedBool, ok := retained.(bool)

			if !ok {
				return nil, fmt.Errorf("mqtt.message.create retained must be a boolean")
			}

			//TODO(jwetzell): convert payload into []byte or string for sending
			payload, ok := params["payload"]

			if !ok {
				return nil, fmt.Errorf("mqtt.message.create requires an payload parameter")
			}

			return &MQTTMessageCreate{config: config, Topic: topicString, QoS: byte(qosByte), Retained: retainedBool, Payload: payload}, nil
		},
	})
}
