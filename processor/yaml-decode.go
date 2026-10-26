package processor

import (
	"context"
	"errors"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/params"
	"github.com/jwetzell/showbridge-go/routing"
	"sigs.k8s.io/yaml"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:        "yaml.decode",
		Title:       "Decode YAML",
		Description: "Decodes incoming bytes as YAML",
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			return &YamlDecode{config: config}, nil
		},
	})
}

type YamlDecode struct {
	config config.ProcessorConfig
}

func (jd *YamlDecode) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
	payload := wrappedPayload.Payload

	payloadBytes, ok := params.GetAnyAsByteSlice(payload)

	if !ok {
		return wrappedPayload, errors.New("yaml.decode can only process a string or []byte")
	}

	payloadYaml := make(map[string]any)

	err := yaml.Unmarshal(payloadBytes, &payloadYaml)
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	wrappedPayload.Payload = payloadYaml
	return wrappedPayload, nil

}
