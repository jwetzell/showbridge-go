package processor

import (
	"context"
	"errors"
	"strconv"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type IntParse struct {
	Base    int
	BitSize int
	config  config.ProcessorConfig
}

func (ip *IntParse) Process(ctx context.Context, payload any) (any, error) {
	payloadString, ok := payload.(string)

	if !ok {
		return nil, errors.New("int.parse processor only accepts a string")
	}

	payloadInt, err := strconv.ParseInt(payloadString, ip.Base, ip.BitSize)
	if err != nil {
		return nil, err
	}
	return payloadInt, nil
}

func (ip *IntParse) Type() string {
	return ip.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "int.parse",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			baseNum := 10
			base, ok := params["base"]
			if ok {
				baseFloat, ok := base.(float64)

				if !ok {
					return nil, errors.New("int.parse base must be a number")
				}

				baseNum = int(baseFloat)
			}

			bitSizeNum := 64
			bitSize, ok := params["bitSize"]
			if ok {
				bitSizeFloat, ok := bitSize.(float64)

				if !ok {
					return nil, errors.New("int.parse bitSize must be a number")
				}

				bitSizeNum = int(bitSizeFloat)
			}
			return &IntParse{config: config, Base: baseNum, BitSize: bitSizeNum}, nil
		},
	})
}
