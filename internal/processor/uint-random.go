package processor

import (
	"context"
	"errors"
	"math/rand/v2"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type UintRandom struct {
	Min    uint
	Max    uint
	config config.ProcessorConfig
}

func (ur *UintRandom) Process(ctx context.Context, payload any) (any, error) {
	payloadInt := rand.UintN(ur.Max-ur.Min+1) + ur.Min
	return payloadInt, nil
}

func (ur *UintRandom) Type() string {
	return ur.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "uint.random",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			min, ok := params["min"]
			if !ok {
				return nil, errors.New("uint.random requires a min parameter")
			}

			minFloat, ok := min.(float64)

			if !ok {
				return nil, errors.New("uint.random min must be a number")
			}

			max, ok := params["max"]
			if !ok {
				return nil, errors.New("uint.random requires a max parameter")
			}

			maxFloat, ok := max.(float64)

			if !ok {
				return nil, errors.New("uint.random max must be a number")
			}

			if maxFloat < minFloat {
				return nil, errors.New("uint.random max must be greater than min")
			}

			return &UintRandom{config: config, Min: uint(minFloat), Max: uint(maxFloat)}, nil
		},
	})
}
