package processor

import (
	"context"
	"errors"
	"fmt"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type IntScale struct {
	OutMin int
	OutMax int
	InMin  int
	InMax  int
	config config.ProcessorConfig
}

func (ir *IntScale) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadInt, ok := common.GetAnyAs[int](payload)
	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("int.scale can only process an int")
	}

	payloadInt = (payloadInt-ir.InMin)*(ir.OutMax-ir.OutMin)/(ir.InMax-ir.InMin) + ir.OutMin
	wrappedPayload.Payload = payloadInt
	return wrappedPayload, nil
}

func (ir *IntScale) Type() string {
	return ir.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "int.scale",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			inMinInt, err := params.GetInt("inMin")
			if err != nil {
				return nil, fmt.Errorf("int.scale inMin error: %w", err)
			}

			inMaxInt, err := params.GetInt("inMax")
			if err != nil {
				return nil, fmt.Errorf("int.scale inMax error: %w", err)
			}

			if inMaxInt < inMinInt {
				return nil, errors.New("int.scale inMax must be greater than inMin")
			}

			outMinInt, err := params.GetInt("outMin")
			if err != nil {
				return nil, fmt.Errorf("int.scale outMin error: %w", err)
			}

			outMaxInt, err := params.GetInt("outMax")
			if err != nil {
				return nil, fmt.Errorf("int.scale outMax error: %w", err)
			}

			if outMaxInt < outMinInt {
				return nil, errors.New("int.scale outMax must be greater than outMin")
			}

			return &IntScale{config: config, InMin: inMinInt, InMax: inMaxInt, OutMin: outMinInt, OutMax: outMaxInt}, nil
		},
	})
}
