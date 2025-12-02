package processing

import (
	"bytes"
	"context"
	"fmt"
	"text/template"
)

type StringCreate struct {
	config   ProcessorConfig
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
		New: func(config ProcessorConfig) (Processor, error) {
			params := config.Params
			tmpl, ok := params["template"]

			if !ok {
				return nil, fmt.Errorf("string.create requires a template parameter")
			}

			templateString, ok := tmpl.(string)

			if !ok {
				return nil, fmt.Errorf("string.create template must be a string")
			}

			templateTemplate, err := template.New("template").Parse(templateString)

			if err != nil {
				return nil, err
			}
			return &StringCreate{config: config, Template: templateTemplate}, nil
		},
	})
}
