package processor

import (
	"context"
	"errors"
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
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params
			bitSizeNum := 64
			bitSize, ok := params["bitSize"]
			if ok {
				bitSizeFloat, ok := bitSize.(float64)

				if !ok {
					return nil, errors.New("float.parse bitSize must be a number")
				}

				bitSizeNum = int(bitSizeFloat)
			}
			return &FloatParse{config: config, BitSize: bitSizeNum}, nil
		},
	})
}
