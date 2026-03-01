package processor

import (
	"context"
	"errors"
	"fmt"
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

			minInt, err := params.GetInt("min")
			if err != nil {
				return nil, fmt.Errorf("uint.random min error: %w", err)
			}

			maxInt, err := params.GetInt("max")
			if err != nil {
				return nil, fmt.Errorf("uint.random max error: %w", err)
			}

			if maxInt < minInt {
				return nil, errors.New("uint.random max must be greater than min")
			}

			return &UintRandom{config: config, Min: uint(minInt), Max: uint(maxInt)}, nil
		},
	})
}
