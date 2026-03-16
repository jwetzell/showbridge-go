package processor

import (
	"context"
	"errors"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type StringDecode struct {
	config config.ProcessorConfig
}

func (sd *StringDecode) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadBytes, ok := common.GetAnyAs[[]byte](payload)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("string.decode processor only accepts a []byte")
	}

	payloadMessage := string(payloadBytes)

	wrappedPayload.Payload = payloadMessage
	return wrappedPayload, nil
}

func (sd *StringDecode) Type() string {
	return sd.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "string.decode",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &StringDecode{config: config}, nil
		},
	})
}
