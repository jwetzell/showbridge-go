package processor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type IntParse struct {
	Base    int
	BitSize int
	config  config.ProcessorConfig
}

func (ip *IntParse) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadString, ok := common.GetAnyAs[string](payload)

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

func (ip *IntParse) Type() string {
	return ip.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "int.parse",
		Title: "Parse Int",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"base": {
					Title:   "Base",
					Type:    "integer",
					Enum:    []any{0, 2, 8, 10, 16},
					Default: json.RawMessage("10"),
				},
				"bitSize": {
					Title:   "Bit Size",
					Type:    "integer",
					Enum:    []any{0, 8, 16, 32, 64},
					Default: json.RawMessage("64"),
				},
			},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(moduleConfig config.ProcessorConfig) (Processor, error) {
			params := moduleConfig.Params

			baseNum, err := params.GetInt("base")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					baseNum = 10
				} else {
					return nil, fmt.Errorf("int.parse base error: %w", err)
				}
			}

			bitSizeNum, err := params.GetInt("bitSize")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					bitSizeNum = 64
				} else {
					return nil, fmt.Errorf("int.parse bitSize error: %w", err)
				}
			}
			return &IntParse{config: moduleConfig, Base: baseNum, BitSize: bitSizeNum}, nil
		},
	})
}
