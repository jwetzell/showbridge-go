package processor

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type StringSplit struct {
	config    config.ProcessorConfig
	Separator string
}

func (ss *StringSplit) Process(ctx context.Context, payload any) (any, error) {
	payloadString, ok := payload.(string)

	if !ok {
		return nil, errors.New("string.split only accepts a string")
	}

	payloadParts := strings.Split(payloadString, ss.Separator)

	return payloadParts, nil
}

func (ss *StringSplit) Type() string {
	return ss.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "string.split",
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
