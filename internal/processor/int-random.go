package processor

import (
	"context"
	"errors"
	"math/rand/v2"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type IntRandom struct {
	Min    int
	Max    int
	config config.ProcessorConfig
}

func (up *IntRandom) Process(ctx context.Context, payload any) (any, error) {
	payloadInt := rand.IntN(up.Max-up.Min+1) + up.Min
	return payloadInt, nil
}

func (up *IntRandom) Type() string {
	return up.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "int.random",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			min, ok := params["min"]
			if !ok {
				return nil, errors.New("int.random requires a min parameter")
			}

			minFloat, ok := min.(float64)

			if !ok {
				return nil, errors.New("int.random min must be a number")
			}

			max, ok := params["max"]
			if !ok {
				return nil, errors.New("int.random requires a max parameter")
			}

			maxFloat, ok := max.(float64)

			if !ok {
				return nil, errors.New("int.random max must be a number")
			}

			if maxFloat < minFloat {
				return nil, errors.New("int.random max must be greater than min")
			}

			return &IntRandom{config: config, Min: int(minFloat), Max: int(maxFloat)}, nil
		},
	})
}
