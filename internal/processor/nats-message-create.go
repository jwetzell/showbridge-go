package processor

import (
	"bytes"
	"context"
	"errors"
	"text/template"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type NATSMessage struct {
	Subject string
	Payload []byte
}

type NATSMessageCreate struct {
	config  config.ProcessorConfig
	Subject *template.Template
	Payload *template.Template
}

func (nmc *NATSMessageCreate) Process(ctx context.Context, payload any) (any, error) {

	var payloadBuffer bytes.Buffer
	err := nmc.Payload.Execute(&payloadBuffer, payload)

	if err != nil {
		return nil, err
	}

	payloadString := payloadBuffer.String()

	var subjectBuffer bytes.Buffer
	err = nmc.Subject.Execute(&subjectBuffer, payload)

	if err != nil {
		return nil, err
	}

	subjectString := subjectBuffer.String()

	message := NATSMessage{
		Subject: subjectString,
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
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params
			subject, ok := params["subject"]

			if !ok {
				return nil, errors.New("nats.message.create requires a subject parameter")
			}

			subjectString, ok := subject.(string)

			if !ok {
				return nil, errors.New("nats.message.create subject must be a string")
			}

			subjectTemplate, err := template.New("subject").Parse(subjectString)

			if err != nil {
				return nil, err
			}

			payload, ok := params["payload"]

			if !ok {
				return nil, errors.New("nats.message.create requires a payload parameter")
			}

			payloadString, ok := payload.(string)

			if !ok {
				return nil, errors.New("nats.message.create payload must be a string")
			}

			payloadTemplate, err := template.New("payload").Parse(payloadString)

			if err != nil {
				return nil, err
			}

			return &NATSMessageCreate{config: config, Subject: subjectTemplate, Payload: payloadTemplate}, nil
		},
	})
}
