package processor

import (
	"context"
	"errors"
	"strconv"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type FloatParse struct {
	config config.ProcessorConfig
}

func (fp *FloatParse) Process(ctx context.Context, payload any) (any, error) {
	payloadString, ok := payload.(string)

	if !ok {
		return nil, errors.New("float.parse processor only accepts a string")
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
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &FloatParse{config: config}, nil
		},
	})
}
