package processor

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
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

func (nmc *NATSMessageCreate) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {

	templateData := wrappedPayload

	var payloadBuffer bytes.Buffer
	err := nmc.Payload.Execute(&payloadBuffer, templateData)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	payloadString := payloadBuffer.String()

	var subjectBuffer bytes.Buffer
	err = nmc.Subject.Execute(&subjectBuffer, templateData)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	subjectString := subjectBuffer.String()

	wrappedPayload.Payload = NATSMessage{
		Subject: subjectString,
		Payload: []byte(payloadString),
	}

	return wrappedPayload, nil
}

func (nmc *NATSMessageCreate) Type() string {
	return nmc.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "nats.message.create",
		Title: "Create NATS Message",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"subject": {
					Title: "Subject",
					Type:  "string",
				},
				"payload": {
					Title: "Payload",
					Type:  "string",
				},
			},
			Required:             []string{"subject", "payload"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params
			subjectString, err := params.GetString("subject")
			if err != nil {
				return nil, fmt.Errorf("nats.message.create subject error: %w", err)
			}

			subjectTemplate, err := template.New("subject").Parse(subjectString)

			if err != nil {
				return nil, err
			}

			payloadString, err := params.GetString("payload")
			if err != nil {
				return nil, fmt.Errorf("nats.message.create payload error: %w", err)
			}

			payloadTemplate, err := template.New("payload").Parse(payloadString)

			if err != nil {
				return nil, err
			}

			return &NATSMessageCreate{config: config, Subject: subjectTemplate, Payload: payloadTemplate}, nil
		},
	})
}
