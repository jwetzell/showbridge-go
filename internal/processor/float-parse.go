package processor

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type FloatParse struct {
	BitSize int
	config  config.ProcessorConfig
}

func (fp *FloatParse) Process(ctx context.Context, payload any) (any, error) {
	payloadString, ok := payload.(string)

	if !ok {
		return nil, errors.New("float.parse processor only accepts a string")
	}

	payloadFloat, err := strconv.ParseFloat(payloadString, fp.BitSize)
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
		New: func(moduleConfig config.ProcessorConfig) (Processor, error) {
			params := moduleConfig.Params

			bitSizeNum, err := params.GetInt("bitSize")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					bitSizeNum = 64
				} else {
					return nil, fmt.Errorf("float.parse bitSize error: %w", err)
				}
			}
			return &FloatParse{config: moduleConfig, BitSize: bitSizeNum}, nil
		},
	})
}
