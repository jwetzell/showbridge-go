package processor

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type IntRandom struct {
	Min    int
	Max    int
	config config.ProcessorConfig
}

func (ir *IntRandom) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	payloadInt := rand.IntN(ir.Max-ir.Min+1) + ir.Min
	wrappedPayload.Payload = payloadInt
	return wrappedPayload, nil
}

func (ir *IntRandom) Type() string {
	return ir.config.Type
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
