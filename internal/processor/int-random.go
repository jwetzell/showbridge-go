package processor

import (
	"context"
	"errors"
	"fmt"
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

			minInt, err := params.GetInt("min")
			if err != nil {
				return nil, fmt.Errorf("int.random min error: %w", err)
			}

			maxInt, err := params.GetInt("max")
			if err != nil {
				return nil, fmt.Errorf("int.random max error: %w", err)
			}

			if maxInt < minInt {
				return nil, errors.New("int.random max must be greater than min")
			}

			return &IntRandom{config: config, Min: int(minInt), Max: int(maxInt)}, nil
		},
	})
}
