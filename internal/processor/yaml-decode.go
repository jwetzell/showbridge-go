package processor

import (
	"context"
	"errors"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/common"
	"sigs.k8s.io/yaml"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "yaml.decode",
		Title: "Decode YAML",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &YamlDecode{config: config}, nil
		},
	})
}

type YamlDecode struct {
	config config.ProcessorConfig
}

func (jd *YamlDecode) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload

	payloadBytes, ok := common.GetAnyAsByteSlice(payload)

	if !ok {
		payloadString, ok := common.GetAnyAs[string](payload)
		if !ok {
			wrappedPayload.End = true
			return wrappedPayload, errors.New("yaml.decode can only process a string or []byte")
		}
		payloadBytes = []byte(payloadString)
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

func (jd *YamlDecode) Id() string {
	return jd.config.Id
}

func (jd *YamlDecode) Type() string {
	return jd.config.Type
}
