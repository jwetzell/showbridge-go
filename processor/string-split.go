package processor

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "string.split",
		Title: "Split String",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"separator": {
					Title:       "Separator",
					Description: "separator to split the string on",
					Type:        "string",
				},
			},
			Required:             []string{"separator"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			processorParams := config.Params

			separatorString, err := processorParams.GetString("separator")
			if err != nil {
				return nil, fmt.Errorf("string.split separator error: %w", err)
			}

			return &StringSplit{config: config, Separator: separatorString}, nil
		},
	})
}

type StringSplit struct {
	config    config.ProcessorConfig
	Separator string
}

func (ss *StringSplit) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadString, ok := payload.(string)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("string.split only accepts a string")
	}

	wrappedPayload.Payload = strings.Split(payloadString, ss.Separator)

	return wrappedPayload, nil
}
