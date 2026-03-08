package processor

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type JsonDecode struct {
	config config.ProcessorConfig
}

func (jd *JsonDecode) Process(ctx context.Context, payload any) (any, error) {

	payloadBytes, ok := GetAnyAsByteSlice(payload)

	if !ok {
		payloadString, ok := GetAnyAs[string](payload)
		if !ok {
			return nil, errors.New("json.decode can only process a string or []byte")
		}
		payloadBytes = []byte(payloadString)
	}

	payloadJson := make(map[string]any)

	err := json.Unmarshal(payloadBytes, &payloadJson)
	if err != nil {
		return nil, err
	}

	return payloadJson, nil

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
