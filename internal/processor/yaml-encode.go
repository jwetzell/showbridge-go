package processor

import (
	"context"

	"sigs.k8s.io/yaml"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/common"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "yaml.encode",
		Title: "Encode YAML",
		New: func(config config.ProcessorConfig) (common.Processor, error) {
			return &YamlEncode{config: config}, nil
		},
	})
}

type YamlEncode struct {
	config config.ProcessorConfig
}

func (je *YamlEncode) Process(ctx context.Context, wrappedPayload common.WrappedPayload, moduleManager common.ModuleManager, inputHandler common.InputHandler) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload

	payloadBytes, err := yaml.Marshal(payload)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	wrappedPayload.Payload = payloadBytes
	return wrappedPayload, nil
}
