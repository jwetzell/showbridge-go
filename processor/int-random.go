package processor

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "int.random",
		Title: "Random Int",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"min": {
					Title:       "Minimum",
					Description: "inclusive minimum value of the random integer",
					Type:        "integer",
				},
				"max": {
					Title:       "Maximum",
					Description: "inclusive maximum value of the random integer",
					Type:        "integer",
				},
			},
			Required:             []string{"min", "max"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			processorParams := config.Params

			minInt, err := processorParams.GetInt("min")
			if err != nil {
				return nil, fmt.Errorf("int.random min error: %w", err)
			}

			maxInt, err := processorParams.GetInt("max")
			if err != nil {
				return nil, fmt.Errorf("int.random max error: %w", err)
			}

			if maxInt < minInt {
				return nil, errors.New("int.random max must be greater than min")
			}

			return &IntRandom{config: config, Min: int(minInt), Max: int(maxInt)}, nil
		},
	})
}

type IntRandom struct {
	Min    int
	Max    int
	config config.ProcessorConfig
}

func (ir *IntRandom) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
	payloadInt := rand.IntN(ir.Max-ir.Min+1) + ir.Min
	wrappedPayload.Payload = payloadInt
	return wrappedPayload, nil
}
