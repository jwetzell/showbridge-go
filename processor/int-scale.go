package processor

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/params"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:        "int.scale",
		Title:       "Scale Int",
		Description: "Scales an integer from one range to another",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"inMin": {
					Title:       "Input Minimum",
					Description: "minimum value of the input integer",
					Type:        "integer",
				},
				"inMax": {
					Title:       "Input Maximum",
					Description: "maximum value of the input integer",
					Type:        "integer",
				},
				"outMin": {
					Title:       "Output Minimum",
					Description: "minimum value of the output integer",
					Type:        "integer",
				},
				"outMax": {
					Title:       "Output Maximum",
					Description: "maximum value of the output integer",
					Type:        "integer",
				},
			},
			Required:             []string{"inMin", "inMax", "outMin", "outMax"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			processorParams := config.Params

			inMinInt, err := processorParams.GetInt("inMin")
			if err != nil {
				return nil, fmt.Errorf("int.scale inMin error: %w", err)
			}

			inMaxInt, err := processorParams.GetInt("inMax")
			if err != nil {
				return nil, fmt.Errorf("int.scale inMax error: %w", err)
			}

			if inMaxInt < inMinInt {
				return nil, errors.New("int.scale inMax must be greater than inMin")
			}

			outMinInt, err := processorParams.GetInt("outMin")
			if err != nil {
				return nil, fmt.Errorf("int.scale outMin error: %w", err)
			}

			outMaxInt, err := processorParams.GetInt("outMax")
			if err != nil {
				return nil, fmt.Errorf("int.scale outMax error: %w", err)
			}

			if outMaxInt < outMinInt {
				return nil, errors.New("int.scale outMax must be greater than outMin")
			}

			return &IntScale{config: config, InMin: inMinInt, InMax: inMaxInt, OutMin: outMinInt, OutMax: outMaxInt}, nil
		},
	})
}

type IntScale struct {
	OutMin int
	OutMax int
	InMin  int
	InMax  int
	config config.ProcessorConfig
}

func (is *IntScale) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadInt, ok := params.GetAnyAsInt(payload)
	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("int.scale can only process an int")
	}

	payloadInt = (payloadInt-is.InMin)*(is.OutMax-is.OutMin)/(is.InMax-is.InMin) + is.OutMin
	wrappedPayload.Payload = payloadInt
	return wrappedPayload, nil
}
