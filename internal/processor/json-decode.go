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
	payloadString, ok := payload.(string)

	if !ok {
		return nil, errors.New("json.decode processor only accepts a string")
	}

	payloadJson := make(map[string]any)

	err := json.Unmarshal([]byte(payloadString), &payloadJson)
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
