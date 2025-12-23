package processor

import (
	"context"
	"errors"
	"strconv"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type IntParse struct {
	config config.ProcessorConfig
}

func (ip *IntParse) Process(ctx context.Context, payload any) (any, error) {
	payloadString, ok := payload.(string)

	if !ok {
		return nil, errors.New("int.parse processor only accepts a string")
	}

	// TODO(jwetzell): make base and bitSize configurable
	payloadInt, err := strconv.ParseInt(payloadString, 10, 64)
	if err != nil {
		return nil, err
	}
	return payloadInt, nil
}

func (ip *IntParse) Type() string {
	return ip.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "int.parse",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &IntParse{config: config}, nil
		},
	})
}
