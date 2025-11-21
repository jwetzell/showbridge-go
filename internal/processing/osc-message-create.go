package processing

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/jwetzell/osc-go"
)

type OSCMessageCreate struct {
	config  ProcessorConfig
	Address *template.Template
}

func (o *OSCMessageCreate) Process(ctx context.Context, payload any) (any, error) {

	var addressBuffer bytes.Buffer
	err := o.Address.Execute(&addressBuffer, payload)

	if err != nil {
		return nil, err
	}

	addressString := addressBuffer.String()

	if len(addressString) == 0 {
		return nil, fmt.Errorf("osc.message.create address must not be empty")
	}

	if addressString[0] != '/' {
		return nil, fmt.Errorf("osc.message.create address must start with '/'")
	}

	payloadMessage := osc.OSCMessage{
		Address: addressString,
	}

	return payloadMessage, nil
}

func (o *OSCMessageCreate) Type() string {
	return o.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "osc.message.create",
		New: func(config ProcessorConfig) (Processor, error) {
			params := config.Params
			address, ok := params["address"]

			if !ok {
				return nil, fmt.Errorf("osc.message.create requires an address parameter")
			}

			addressString, ok := address.(string)

			if !ok {
				return nil, fmt.Errorf("osc.message.create address must be a string")
			}

			addressTemplate, err := template.New("address").Parse(addressString)

			if err != nil {
				return nil, err
			}

			return &OSCMessageCreate{config: config, Address: addressTemplate}, nil
		},
	})
}
