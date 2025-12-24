package processor

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"text/template"

	"github.com/jwetzell/osc-go"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type OSCMessageCreate struct {
	config  config.ProcessorConfig
	Address *template.Template
	Args    []*template.Template
	Types   string
}

func (omc *OSCMessageCreate) Process(ctx context.Context, payload any) (any, error) {

	var addressBuffer bytes.Buffer
	err := omc.Address.Execute(&addressBuffer, payload)

	if err != nil {
		return nil, err
	}

	addressString := addressBuffer.String()

	if len(addressString) == 0 {
		return nil, errors.New("osc.message.create address must not be empty")
	}

	if addressString[0] != '/' {
		return nil, errors.New("osc.message.create address must start with '/'")
	}

	payloadMessage := osc.OSCMessage{
		Address: addressString,
	}

	args := []osc.OSCArg{}

	for argIndex, argTemplate := range omc.Args {
		var argBuffer bytes.Buffer
		err := argTemplate.Execute(&argBuffer, payload)

		if err != nil {
			return nil, err
		}

		argString := argBuffer.String()

		typedArg, err := argToTypedArg(argString, omc.Types[argIndex])

		if err != nil {
			return nil, err
		}

		args = append(args, typedArg)
	}

	if len(args) > 0 {
		payloadMessage.Args = args
	}

	return payloadMessage, nil
}

func (omc *OSCMessageCreate) Type() string {
	return omc.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "osc.message.create",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params
			address, ok := params["address"]

			if !ok {
				return nil, errors.New("osc.message.create requires an address parameter")
			}

			addressString, ok := address.(string)

			if !ok {
				return nil, errors.New("osc.message.create address must be a string")
			}

			addressTemplate, err := template.New("address").Parse(addressString)

			if err != nil {
				return nil, err
			}

			args, ok := params["args"]

			if ok {
				rawArgs, ok := args.([]interface{})

				if !ok {
					return nil, fmt.Errorf("osc.message.create address must be an array found %T", args)
				}

				types, ok := params["types"]

				if !ok {
					return nil, errors.New("osc.message.create requires a types parameter with args")
				}

				typesString, ok := types.(string)

				if !ok {
					return nil, errors.New("osc.message.create types must be a string")
				}

				if len(rawArgs) != len(typesString) {
					return nil, errors.New("osc.message.create args and types must be the same length")
				}

				argTemplates := []*template.Template{}

				for _, rawArg := range rawArgs {
					argString, ok := rawArg.(string)

					if !ok {
						return nil, errors.New("osc.message.create arg must be a string")
					}

					argTemplate, err := template.New("arg").Parse(argString)

					if err != nil {
						return nil, err
					}
					argTemplates = append(argTemplates, argTemplate)
				}
				return &OSCMessageCreate{config: config, Address: addressTemplate, Args: argTemplates, Types: typesString}, nil
			}
			return &OSCMessageCreate{config: config, Address: addressTemplate}, nil
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
