package processor

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type FloatRandom struct {
	BitSize int
	Min     float64
	Max     float64
	config  config.ProcessorConfig
}

func (fr *FloatRandom) Process(ctx context.Context, payload any) (any, error) {
	if fr.BitSize == 32 {
		payloadFloat := rand.Float32()*(float32(fr.Max)-float32(fr.Min)) + float32(fr.Min)
		return payloadFloat, nil
	}
	if fr.BitSize == 64 {
		payloadFloat := rand.Float64()*(fr.Max-fr.Min) + fr.Min
		return payloadFloat, nil
	}
	return nil, errors.New("float.random bitSize error: must be 32 or 64")
}

func (fr *FloatRandom) Type() string {
	return fr.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "float.random",
		New: func(processorConfig config.ProcessorConfig) (Processor, error) {
			params := processorConfig.Params

			bitSizeInt, err := params.GetInt("bitSize")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					bitSizeInt = 32
				} else {
					return nil, fmt.Errorf("float.random bitSize error: %w", err)
				}
			}

			if bitSizeInt != 32 && bitSizeInt != 64 {
				return nil, errors.New("float.random bitSize error: must be 32 or 64")
			}

			minFloat, err := params.GetFloat64("min")
			if err != nil {
				return nil, fmt.Errorf("float.random min error: %w", err)
			}

			maxFloat, err := params.GetFloat64("max")
			if err != nil {
				return nil, fmt.Errorf("float.random max error: %w", err)
			}

			if maxFloat < minFloat {
				return nil, errors.New("float.random max must be greater than min")
			}

			return &FloatRandom{config: processorConfig, Min: minFloat, Max: maxFloat, BitSize: bitSizeInt}, nil
		},
	})
}
