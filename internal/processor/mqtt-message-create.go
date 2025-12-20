package processor

import (
	"context"
	"fmt"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type MQTTMessage struct {
	topic    string
	qos      byte
	payload  []byte
	retained bool
}

type MQTTMessageCreate struct {
	config   config.ProcessorConfig
	Topic    string
	QoS      byte
	Retained bool
	Payload  []byte
}

func (mm MQTTMessage) Duplicate() bool {
	// TODO(jwetzell): implement?
	return false
}

func (mm MQTTMessage) Qos() byte {
	return mm.qos
}

func (mm MQTTMessage) Retained() bool {
	return mm.retained
}

func (mm MQTTMessage) Topic() string {
	return mm.topic
}

func (mm MQTTMessage) MessageID() uint16 {
	// TODO(jwetzell): implement?
	return 0
}

func (mm MQTTMessage) Payload() []byte {
	return mm.payload
}

func (mm MQTTMessage) Ack() {}

func (mmc *MQTTMessageCreate) Process(ctx context.Context, payload any) (any, error) {

	message := MQTTMessage{
		topic:    mmc.Topic,
		qos:      mmc.QoS,
		retained: mmc.Retained,
		payload:  mmc.Payload,
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
				return nil, fmt.Errorf("mqtt.message.create requires a topic parameter")
			}

			topicString, ok := topic.(string)

			if !ok {
				return nil, fmt.Errorf("mqtt.message.create topic must be a string")
			}

			qos, ok := params["qos"]

			if !ok {
				return nil, fmt.Errorf("mqtt.message.create requires a qos parameter")
			}

			qosByte, ok := qos.(float64)

			if !ok {
				return nil, fmt.Errorf("mqtt.message.create qos must be a number")
			}

			retained, ok := params["retained"]

			if !ok {
				return nil, fmt.Errorf("mqtt.message.create requires a retained parameter")
			}

			retainedBool, ok := retained.(bool)

			if !ok {
				return nil, fmt.Errorf("mqtt.message.create retained must be a boolean")
			}

			//TODO(jwetzell): convert payload into []byte or string for sending
			payload, ok := params["payload"]

			if !ok {
				return nil, fmt.Errorf("mqtt.message.create requires a payload parameter")
			}

			if payloadBytes, ok := payload.([]byte); ok {
				return &MQTTMessageCreate{config: config, Topic: topicString, QoS: byte(qosByte), Retained: retainedBool, Payload: payloadBytes}, nil
			}

			payloadString, ok := payload.(string)

			if !ok {
				return nil, fmt.Errorf("mqtt.message.create payload must be a string or byte array")
			}

			payloadBytes := []byte(payloadString)

			return &MQTTMessageCreate{config: config, Topic: topicString, QoS: byte(qosByte), Retained: retainedBool, Payload: payloadBytes}, nil
		},
	})
}
