package processor

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type JsonDecode struct {
	config config.ProcessorConfig
}

func (jd *JsonDecode) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload

	payloadBytes, ok := common.GetAnyAsByteSlice(payload)

	if !ok {
		payloadString, ok := common.GetAnyAs[string](payload)
		if !ok {
			wrappedPayload.End = true
			return wrappedPayload, errors.New("json.decode can only process a string or []byte")
		}
		payloadBytes = []byte(payloadString)
	}

	payloadJson := make(map[string]any)

	err := json.Unmarshal(payloadBytes, &payloadJson)
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	wrappedPayload.Payload = payloadJson
	return wrappedPayload, nil

}

func (jd *JsonDecode) Type() string {
	return jd.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "json.decode",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &JsonDecode{config: config}, nil
		},
	})
}
