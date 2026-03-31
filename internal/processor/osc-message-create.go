package processor

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"text/template"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/osc-go"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type OSCMessageCreate struct {
	config  config.ProcessorConfig
	Address *template.Template
	Args    []*template.Template
	Types   string
}

func (omc *OSCMessageCreate) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {

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

	payloadMessage := &osc.OSCMessage{
		Address: addressString,
	}

	args := []osc.OSCArg{}

	for argIndex, argTemplate := range omc.Args {
		var argBuffer bytes.Buffer
		err := argTemplate.Execute(&argBuffer, templateData)

		if err != nil {
			wrappedPayload.End = true
			return wrappedPayload, err
		}

		argString := argBuffer.String()

		typedArg, err := argToTypedArg(argString, omc.Types[argIndex])

		if err != nil {
			wrappedPayload.End = true
			return wrappedPayload, err
		}

		args = append(args, typedArg)
	}

	if len(args) > 0 {
		payloadMessage.Args = args
	}

	wrappedPayload.Payload = payloadMessage
	return wrappedPayload, nil
}

func (omc *OSCMessageCreate) Type() string {
	return omc.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "osc.message.create",
		Title: "Create OSC Message",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"address": {
					Title: "Address",
					Type:  "string",
				},
				"args": {
					Title: "Arguments",
					Type:  "array",
					Items: &jsonschema.Schema{
						Type: "string",
					},
				},
				"types": {
					Title: "Argument Types",
					Type:  "string",
				},
			},
			Required:             []string{"address"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(processorConfig config.ProcessorConfig) (Processor, error) {
			params := processorConfig.Params
			addressString, err := params.GetString("address")
			if err != nil {
				return nil, fmt.Errorf("osc.message.create address error: %w", err)
			}

			addressTemplate, err := template.New("address").Parse(addressString)

			if err != nil {
				return nil, err
			}

			argStrings, err := params.GetStringSlice("args")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					return &OSCMessageCreate{config: processorConfig, Address: addressTemplate}, nil
				} else {
					return nil, fmt.Errorf("osc.message.create args error: %w", err)
				}
			}

			typesString, err := params.GetString("types")
			if err != nil {
				return nil, fmt.Errorf("osc.message.create types error: %w", err)
			}

			if len(argStrings) != len(typesString) {
				return nil, errors.New("osc.message.create args and types must be the same length")
			}

			argTemplates := []*template.Template{}

			for _, argString := range argStrings {

				argTemplate, err := template.New("arg").Parse(argString)

				if err != nil {
					return nil, err
				}
				argTemplates = append(argTemplates, argTemplate)
			}
			return &OSCMessageCreate{config: processorConfig, Address: addressTemplate, Args: argTemplates, Types: typesString}, nil
		},
	})
}

func argToTypedArg(rawArg string, oscType byte) (osc.OSCArg, error) {

	switch oscType {
	case 's':
		return osc.OSCArg{
			Value: rawArg,
			Type:  "s",
		}, nil
	case 'i':
		number, err := strconv.ParseInt(rawArg, 10, 32)
		if err != nil {
			return osc.OSCArg{}, err
		}
		return osc.OSCArg{
			Value: int32(number),
			Type:  "i",
		}, nil
	case 'f':
		number, err := strconv.ParseFloat(rawArg, 32)
		if err != nil {
			return osc.OSCArg{}, err
		}
		return osc.OSCArg{
			Value: float32(number),
			Type:  "f",
		}, nil
	case 'b':
		data, err := hex.DecodeString(rawArg)
		if err != nil {
			return osc.OSCArg{}, err
		}
		return osc.OSCArg{
			Value: data,
			Type:  "b",
		}, nil
	case 'h':
		number, err := strconv.ParseInt(rawArg, 10, 64)
		if err != nil {
			return osc.OSCArg{}, err
		}
		return osc.OSCArg{
			Value: int64(number),
			Type:  "h",
		}, nil
	case 'd':
		number, err := strconv.ParseFloat(rawArg, 64)
		if err != nil {
			return osc.OSCArg{}, err
		}
		return osc.OSCArg{
			Value: float64(number),
			Type:  "d",
		}, nil
	case 'T':
		return osc.OSCArg{
			Value: true,
			Type:  "T",
		}, nil
	case 'F':
		return osc.OSCArg{
			Value: false,
			Type:  "F",
		}, nil
	case 'N':
		return osc.OSCArg{
			Value: nil,
			Type:  "N",
		}, nil
	default:
		return osc.OSCArg{}, fmt.Errorf("osc.message.create unhandled osc type: %c", oscType)
	}
}
