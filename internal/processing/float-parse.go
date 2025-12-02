package processing

import (
	"context"
	"fmt"
	"strconv"
)

type FloatParse struct {
	config ProcessorConfig
}

func (fp *FloatParse) Process(ctx context.Context, payload any) (any, error) {
	payloadString, ok := payload.(string)

	if !ok {
		return nil, fmt.Errorf("float.parse processor only accepts a string")
	}

	// TODO(jwetzell): make bitSize configurable
	payloadFloat, err := strconv.ParseFloat(payloadString, 64)
	if err != nil {
		return nil, err
	}
	return payloadFloat, nil
}

func (fp *FloatParse) Type() string {
	return fp.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "float.parse",
		New: func(config ProcessorConfig) (Processor, error) {
			return &FloatParse{config: config}, nil
		},
	})
}
