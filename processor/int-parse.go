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
		Type:  "int.parse",
		Title: "Parse Int",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"base": {
					Title:       "Base",
					Description: "numerical base to use for parsing the integer",
					Type:        "integer",
					Enum:        []any{0, 2, 8, 10, 16},
					Default:     json.RawMessage("10"),
				},
				"bitSize": {
					Title:       "Bit Size",
					Description: "bit size of the resulting integer",
					Type:        "integer",
					Enum:        []any{0, 8, 16, 32, 64},
					Default:     json.RawMessage("64"),
				},
			},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(processorConfig config.ProcessorConfig) (routing.Processor, error) {
			processorParams := processorConfig.Params

			baseNum, err := processorParams.GetInt("base")
			if err != nil {
				if errors.Is(err, params.ErrParamNotFound) {
					baseNum = 10
				} else {
					return nil, fmt.Errorf("int.parse base error: %w", err)
				}
			}

			bitSizeNum, err := processorParams.GetInt("bitSize")
			if err != nil {
				if errors.Is(err, params.ErrParamNotFound) {
					bitSizeNum = 64
				} else {
					return nil, fmt.Errorf("int.parse bitSize error: %w", err)
				}
			}
			return &IntParse{config: processorConfig, Base: baseNum, BitSize: bitSizeNum}, nil
		},
	})
}

type IntParse struct {
	Base    int
	BitSize int
	config  config.ProcessorConfig
}

func (ip *IntParse) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, inputHandler routing.InputHandler) (routing.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadString, ok := payload.(string)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("int.parse processor only accepts a string")
	}

	payloadInt, err := strconv.ParseInt(payloadString, ip.Base, ip.BitSize)
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}
	wrappedPayload.Payload = payloadInt
	return wrappedPayload, nil
}
