package processor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type KVGet struct {
	config   config.ProcessorConfig
	ModuleId string
	Key      string
	logger   *slog.Logger
}

func (kvg *KVGet) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	if wrappedPayload.Modules == nil {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("kv.get wrapped payload has no modules")
	}

	module, ok := wrappedPayload.Modules[kvg.ModuleId]
	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("kv.get unable to find module with id: %s", kvg.ModuleId)
	}

	kvModule, ok := module.(common.KeyValueModule)
	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("kv.get module with id %s is not a KeyValueModule", kvg.ModuleId)
	}

	value, err := kvModule.Get(kvg.Key)
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("kv.get error getting key: %w", err)
	}

	wrappedPayload.Payload = value
	return wrappedPayload, nil
}

func (kvg *KVGet) Type() string {
	return kvg.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "kv.get",
		Title: "Get Key",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"module": {
					Title: "Module ID",
					Type:  "string",
				},
				"key": {
					Title: "Key",
					Type:  "string",
				},
			},
			Required:             []string{"module", "key"},
			AdditionalProperties: nil,
		},
		New: func(config config.ProcessorConfig) (Processor, error) {

			params := config.Params

			moduleIdString, err := params.GetString("module")
			if err != nil {
				return nil, fmt.Errorf("kv.get module error: %w", err)
			}

			keyString, err := params.GetString("key")
			if err != nil {
				return nil, fmt.Errorf("kv.get key error: %w", err)
			}
			return &KVGet{config: config, ModuleId: moduleIdString, Key: keyString, logger: slog.Default().With("component", "processor", "type", config.Type)}, nil
		},
	})
}
