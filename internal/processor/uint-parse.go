package processor

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type UintParse struct {
	config config.ProcessorConfig
}

func (up *UintParse) Process(ctx context.Context, payload any) (any, error) {
	payloadString, ok := payload.(string)

	if !ok {
		return nil, fmt.Errorf("uint.parse processor only accepts a string")
	}

	// TODO(jwetzell): make base and bitSize configurable
	payloadUint, err := strconv.ParseUint(payloadString, 10, 64)
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
			return &UintParse{config: config}, nil
		},
	})
}
