package processor

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type StringCreate struct {
	config   config.ProcessorConfig
	Template *template.Template
}

func (sc *StringCreate) Process(ctx context.Context, payload any) (any, error) {
	var templateBuffer bytes.Buffer
	err := sc.Template.Execute(&templateBuffer, payload)

	if err != nil {
		return nil, err
	}

	payloadString := templateBuffer.String()

	return payloadString, nil
}

func (sc *StringCreate) Type() string {
	return sc.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "string.create",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params
			templateString, err := params.GetString("template")
			if err != nil {
				return nil, fmt.Errorf("string.create template error: %w", err)
			}

			templateTemplate, err := template.New("template").Parse(templateString)

			if err != nil {
				return nil, err
			}
			return &StringCreate{config: config, Template: templateTemplate}, nil
		},
	})
}
