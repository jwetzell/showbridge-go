package processor

import (
	"context"
	"errors"

	freeD "github.com/jwetzell/free-d-go"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type FreeDEncode struct {
	config config.ProcessorConfig
}

func (fe *FreeDEncode) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadPosition, ok := common.GetAnyAs[freeD.FreeDPosition](payload)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("freed.decode processor only accepts a FreeDEncode")
	}

	payloadBytes := freeD.Encode(payloadPosition)

	wrappedPayload.Payload = payloadBytes
	return wrappedPayload, nil
}

func (fe *FreeDEncode) Type() string {
	return fe.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "freed.encode",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &FreeDEncode{config: config}, nil
		},
	})
}
