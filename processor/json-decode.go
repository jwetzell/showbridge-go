package processor

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/params"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:        "json.decode",
		Title:       "Decode JSON",
		Description: "Decodes bytes into a JSON object",
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			return &JsonDecode{config: config}, nil
		},
	})
}

type JsonDecode struct {
	config config.ProcessorConfig
}

func (jd *JsonDecode) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
	payload := wrappedPayload.Payload

	payloadBytes, ok := params.GetAnyAsByteSlice(payload)

	if !ok {
		return wrappedPayload, errors.New("json.decode can only process a string or []byte")
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
