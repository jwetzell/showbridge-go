package processor

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type FloatParse struct {
	BitSize int
	config  config.ProcessorConfig
}

func (fp *FloatParse) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadString, ok := common.GetAnyAs[string](payload)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("float.parse processor only accepts a string")
	}

	payloadFloat, err := strconv.ParseFloat(payloadString, fp.BitSize)
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}
	wrappedPayload.Payload = payloadFloat
	return wrappedPayload, nil
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
