package processor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/params"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "float.parse",
		Title: "Parse Float",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"bitSize": {
					Title:       "Bit Size",
					Description: "bit size of the resulting float",
					Type:        "integer",
					Enum:        []any{32, 64},
					Default:     json.RawMessage("64"),
				},
			},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(processorConfig config.ProcessorConfig) (routing.Processor, error) {
			processorParams := processorConfig.Params

			bitSizeNum, err := processorParams.GetInt("bitSize")
			if err != nil {
				if errors.Is(err, params.ErrParamNotFound) {
					bitSizeNum = 64
				} else {
					return nil, fmt.Errorf("float.parse bitSize error: %w", err)
				}
			}
			return &FloatParse{config: processorConfig, BitSize: bitSizeNum}, nil
		},
	})
}

type FloatParse struct {
	BitSize int
	config  config.ProcessorConfig
}

func (fp *FloatParse) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadString, ok := payload.(string)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("float.parse processor only accepts a string")
	}

	payloadFloat, err := strconv.ParseFloat(payloadString, fp.BitSize)
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}
	wrappedPayload.Payload = payloadFloat
	return wrappedPayload, nil
}
