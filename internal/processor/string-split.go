package processor

import (
	"context"
	"errors"
	"strings"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type StringSplit struct {
	config    config.ProcessorConfig
	Separator string
}

func (se *StringSplit) Process(ctx context.Context, payload any) (any, error) {
	payloadString, ok := payload.(string)

	if !ok {
		return nil, errors.New("string.split only accepts a string")
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
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			separator, ok := params["separator"]

			if !ok {
				return nil, errors.New("string.split requires a separator")
			}

			separatorString, ok := separator.(string)

			if !ok {
				return nil, errors.New("string.split separator must be a string")
			}

			return &StringSplit{config: config, Separator: separatorString}, nil
		},
	})
}
