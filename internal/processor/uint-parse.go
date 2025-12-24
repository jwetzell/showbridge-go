package processor

import (
	"context"
	"errors"
	"strconv"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type UintParse struct {
	Base    int
	BitSize int
	config  config.ProcessorConfig
}

func (up *UintParse) Process(ctx context.Context, payload any) (any, error) {
	payloadString, ok := payload.(string)

	if !ok {
		return nil, errors.New("uint.parse processor only accepts a string")
	}

	payloadUint, err := strconv.ParseUint(payloadString, up.Base, up.BitSize)
	if err != nil {
		return nil, err
	}
	return payloadUint, nil
}

func (up *UintParse) Type() string {
	return up.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "uint.parse",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params
			baseNum := 10
			base, ok := params["base"]
			if ok {
				baseFloat, ok := base.(float64)

				if !ok {
					return nil, errors.New("uint.parse base must be a number")
				}

				baseNum = int(baseFloat)
			}

			bitSizeNum := 64
			bitSize, ok := params["bitSize"]
			if ok {
				bitSizeFloat, ok := bitSize.(float64)

				if !ok {
					return nil, errors.New("uint.parse bitSize must be a number")
				}

				bitSizeNum = int(bitSizeFloat)
			}
			return &UintParse{config: config, Base: baseNum, BitSize: bitSizeNum}, nil
		},
	})
}
