package processing

import (
	"bytes"
	"context"
	"fmt"
	"text/template"
)

type NATSMessage struct {
	Subject string
	Payload []byte
}

type NATSMessageCreate struct {
	config  ProcessorConfig
	Subject string
	Payload *template.Template
}

func (nmc *NATSMessageCreate) Process(ctx context.Context, payload any) (any, error) {

	var payloadBuffer bytes.Buffer
	err := nmc.Payload.Execute(&payloadBuffer, payload)

	if err != nil {
		return nil, err
	}

	payloadString := payloadBuffer.String()

	message := NATSMessage{
		Subject: nmc.Subject,
		Payload: []byte(payloadString),
	}

	return message, nil
}

func (nmc *NATSMessageCreate) Type() string {
	return nmc.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "nats.message.create",
		New: func(config ProcessorConfig) (Processor, error) {
			params := config.Params
			// TODO(jwetzell): support template for subject
			subject, ok := params["subject"]

			if !ok {
				return nil, fmt.Errorf("nats.message.create requires a subject parameter")
			}

			subjectString, ok := subject.(string)

			if !ok {
				return nil, fmt.Errorf("nats.message.create subject must be a string")
			}

			payload, ok := params["payload"]

			if !ok {
				return nil, fmt.Errorf("osc.message.create requires a payload parameter")
			}

			payloadString, ok := payload.(string)

			if !ok {
				return nil, fmt.Errorf("osc.message.create payload must be a string")
			}

			payloadTemplate, err := template.New("payload").Parse(payloadString)

			if err != nil {
				return nil, err
			}

			return &NATSMessageCreate{config: config, Subject: subjectString, Payload: payloadTemplate}, nil
		},
	})
}
