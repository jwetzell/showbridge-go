package processor

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type StringSplit struct {
	config    config.ProcessorConfig
	Separator string
}

func (ss *StringSplit) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadString, ok := common.GetAnyAs[string](payload)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("string.split only accepts a string")
	}

	wrappedPayload.Payload = strings.Split(payloadString, ss.Separator)

	return wrappedPayload, nil
}

func (ss *StringSplit) Type() string {
	return ss.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "string.split",
		Title: "Split String",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"separator": {
					Title: "Separator",
					Type:  "string",
				},
			},
			Required:             []string{"separator"},
			AdditionalProperties: nil,
		},
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			separatorString, err := params.GetString("separator")
			if err != nil {
				return nil, fmt.Errorf("string.split separator error: %w", err)
			}

			return &StringSplit{config: config, Separator: separatorString}, nil
		},
	})
}
