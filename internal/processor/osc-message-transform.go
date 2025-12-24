package processor

import (
	"bytes"
	"context"
	"errors"
	"text/template"

	osc "github.com/jwetzell/osc-go"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type OSCMessageTransform struct {
	config  config.ProcessorConfig
	Address *template.Template
}

func (o *OSCMessageTransform) Process(ctx context.Context, payload any) (any, error) {
	payloadMessage, ok := payload.(osc.OSCMessage)

	if !ok {
		return nil, errors.New("osc.message.transform processor only accepts an OSCMessage")
	}

	var addressBuffer bytes.Buffer
	err := o.Address.Execute(&addressBuffer, payloadMessage)

	if err != nil {
		return nil, err
	}

	addressString := addressBuffer.String()

	if len(addressString) == 0 {
		return nil, errors.New("osc.message.transform address must not be empty")
	}

	if addressString[0] != '/' {
		return nil, errors.New("osc.message.transform address must start with '/'")
	}

	payloadMessage.Address = addressString

	return payloadMessage, nil
}

func (o *OSCMessageTransform) Type() string {
	return o.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "osc.message.transform",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params
			address, ok := params["address"]

			if !ok {
				return nil, errors.New("osc.message.transform requires an address parameter")
			}

			addressString, ok := address.(string)

			if !ok {
				return nil, errors.New("osc.message.transform address must be a string")
			}

			addressTemplate, err := template.New("address").Parse(addressString)

			if err != nil {
				return nil, err
			}

			return &OSCMessageTransform{config: config, Address: addressTemplate}, nil
		},
	})
}
