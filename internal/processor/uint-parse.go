package processor

import (
	"context"
	"errors"
	"fmt"
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
		New: func(moduleConfig config.ProcessorConfig) (Processor, error) {
			params := moduleConfig.Params
			baseNum, err := params.GetInt("base")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					baseNum = 10
				} else {
					return nil, fmt.Errorf("uint.parse base error: %w", err)
				}
			}

			bitSizeNum, err := params.GetInt("bitSize")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					bitSizeNum = 64
				} else {
					return nil, fmt.Errorf("uint.parse bitSize error: %w", err)
				}
			}
			return &UintParse{config: moduleConfig, Base: baseNum, BitSize: bitSizeNum}, nil
		},
	})
}
