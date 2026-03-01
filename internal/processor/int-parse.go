package processor

import (
	"context"
	"errors"
	"fmt"
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
		New: func(moduleConfig config.ProcessorConfig) (Processor, error) {
			params := moduleConfig.Params

			baseNum, err := params.GetInt("base")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					baseNum = 10
				} else {
					return nil, fmt.Errorf("int.parse base error: %w", err)
				}
			}

			bitSizeNum, err := params.GetInt("bitSize")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					bitSizeNum = 64
				} else {
					return nil, fmt.Errorf("int.parse bitSize error: %w", err)
				}
			}
			return &IntParse{config: moduleConfig, Base: baseNum, BitSize: bitSizeNum}, nil
		},
	})
}
