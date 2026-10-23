package processor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"text/template"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/osc-go"
	"github.com/jwetzell/showbridge-go/common"
	"github.com/jwetzell/showbridge-go/config"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "osc.message.create",
		Title: "Create OSC Message",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"address": {
					Title:       "Address",
					Description: "OSC address for the message",
					Type:        "string",
				},
				"args": {
					Title:       "Arguments",
					Description: "arguments for the OSC message",
					Type:        "array",
					Items: &jsonschema.Schema{
						Type: "object",
						Properties: map[string]*jsonschema.Schema{
							"value": {
								Title:       "Value",
								Description: "Value of the argument",
								Type:        "string",
							},
							"type": {
								Title:       "Type",
								Description: "OSC type of the argument",
								Type:        "string",
								Enum:        []any{"s", "i", "f", "b", "T", "F", "N", "I", "h", "d", "t"},
								MinLength:   new(1),
								MaxLength:   new(1),
							},
						},
						Required:             []string{"value", "type"},
						AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
					},
				},
			},
			Required:             []string{"address"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(processorConfig config.ProcessorConfig) (common.Processor, error) {
			params := processorConfig.Params
			addressString, err := params.GetString("address")
			if err != nil {
				return nil, fmt.Errorf("osc.message.create address error: %w", err)
			}

			addressTemplate, err := template.New("address").Parse(addressString)

			if err != nil {
				return nil, err
			}

			argObjects, err := params.GetMapSlice("args")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					return &MessageCreate{config: processorConfig, Address: addressTemplate}, nil
				} else {
					return nil, fmt.Errorf("osc.message.create args error: %w", err)
				}
			}

			types := []rune{}
			for _, argObject := range argObjects {
				argType, ok := argObject["type"]
				if !ok {
					return nil, errors.New("osc.message.create arg type error: not found")
				}
				argTypeStr, ok := argType.(string)
				if !ok {
					return nil, errors.New("osc.message.create arg type error: not a string")
				}

				if len(argTypeStr) != 1 {
					return nil, errors.New("osc.message.create arg type error: must be a single character")
				}
				types = append(types, rune(argTypeStr[0]))
			}

			argTemplates := []*template.Template{}

			for _, argObject := range argObjects {
				argValue, ok := argObject["value"]
				if !ok {
					return nil, errors.New("osc.message.create arg value error: not found")
				}
				argValueStr, ok := argValue.(string)
				if !ok {
					return nil, errors.New("osc.message.create arg value error: not a string")
				}

				argTemplate, err := template.New("arg").Parse(argValueStr)

				if err != nil {
					return nil, err
				}
				argTemplates = append(argTemplates, argTemplate)
			}
			return &MessageCreate{config: processorConfig, Address: addressTemplate, Args: argTemplates, Types: types}, nil
		},
	})
}

type MessageCreate struct {
	config  config.ProcessorConfig
	Address *template.Template
	Args    []*template.Template
	Types   []rune
}

func (omc *MessageCreate) Process(ctx context.Context, wrappedPayload common.WrappedPayload, moduleManager common.ModuleManager, inputHandler common.InputHandler) (common.WrappedPayload, error) {

	templateData := wrappedPayload

	var addressBuffer bytes.Buffer
	err := omc.Address.Execute(&addressBuffer, templateData)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	addressString := addressBuffer.String()

	if len(addressString) == 0 {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("osc.message.create address must not be empty")
	}

	if addressString[0] != '/' {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("osc.message.create address must start with '/'")
	}

	payloadMessage := &osc.Message{
		Address: addressString,
	}

	args := []osc.Arg{}

	for argIndex, argTemplate := range omc.Args {
		var argBuffer bytes.Buffer
		err := argTemplate.Execute(&argBuffer, templateData)

		if err != nil {
			wrappedPayload.End = true
			return wrappedPayload, err
		}

		argString := argBuffer.String()

		typedArg, err := osc.ArgFromStringAndType(argString, omc.Types[argIndex])

		if err != nil {
			wrappedPayload.End = true
			return wrappedPayload, fmt.Errorf("osc.message.create arg conversion error: %w", err)
		}

		args = append(args, typedArg)
	}

	if len(args) > 0 {
		payloadMessage.Args = args
	}

	wrappedPayload.Payload = payloadMessage
	return wrappedPayload, nil
}
