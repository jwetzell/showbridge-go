package processor

import (
	"context"

	"sigs.k8s.io/yaml"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:        "yaml.encode",
		Title:       "Encode YAML",
		Description: "Encodes the payload as YAML",
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			return &YamlEncode{config: config}, nil
		},
	})
}

type YamlEncode struct {
	config config.ProcessorConfig
}

func (ye *YamlEncode) Process(next routing.ProcessFunc) routing.ProcessFunc {
	return func(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
		payload := wrappedPayload.Payload

		payloadBytes, err := yaml.Marshal(payload)

		if err != nil {
			return wrappedPayload, err
		}

		wrappedPayload.Payload = payloadBytes
		return next(ctx, wrappedPayload, moduleManager, router)
	}
}
