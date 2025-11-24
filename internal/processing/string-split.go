package processing

import (
	"context"
	"fmt"
	"strings"
)

type StringSplit struct {
	config    ProcessorConfig
	Separator string
}

func (se *StringSplit) Process(ctx context.Context, payload any) (any, error) {
	payloadString, ok := payload.(string)

	if !ok {
		return nil, fmt.Errorf("string.split only accepts a string")
	}

	payloadParts := strings.Split(payloadString, se.Separator)

	return payloadParts, nil
}

func (se *StringSplit) Type() string {
	return se.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "string.split",
		New: func(config ProcessorConfig) (Processor, error) {
			params := config.Params

			separator, ok := params["separator"]

			if !ok {
				return nil, fmt.Errorf("string.split requires a separator")
			}

			separatorString, ok := separator.(string)

			if !ok {
				return nil, fmt.Errorf("string.split separator must be a string")
			}

			return &StringSplit{config: config, Separator: separatorString}, nil
		},
	})
}
