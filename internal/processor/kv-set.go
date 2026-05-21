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

type KVSet struct {
	config   config.ProcessorConfig
	ModuleId string
	Key      string
	logger   *slog.Logger
	module   common.KeyValueModule
}

func (kvs *KVSet) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	if kvs.module == nil {
		if wrappedPayload.Modules == nil {
			wrappedPayload.End = true
			return wrappedPayload, errors.New("kv.set wrapped payload has no modules")
		}

		module, ok := wrappedPayload.Modules[kvs.ModuleId]
		if !ok {
			wrappedPayload.End = true
			return wrappedPayload, fmt.Errorf("kv.set unable to find module with id: %s", kvs.ModuleId)
		}

		kvModule, ok := module.(common.KeyValueModule)
		if !ok {
			wrappedPayload.End = true
			return wrappedPayload, fmt.Errorf("kv.set module with id %s is not a KeyValueModule", kvs.ModuleId)
		}
		kvs.module = kvModule
	}

	err := kvs.module.Set(kvs.Key, wrappedPayload.Payload)
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("kv.set error setting key: %w", err)
	}

	return wrappedPayload, nil
}

func (kvs *KVSet) Type() string {
	return kvs.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "kv.set",
		Title: "Set Key",
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
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ProcessorConfig) (Processor, error) {

			params := config.Params

			moduleIdString, err := params.GetString("module")
			if err != nil {
				return nil, fmt.Errorf("kv.set module error: %w", err)
			}

			keyString, err := params.GetString("key")
			if err != nil {
				return nil, fmt.Errorf("kv.set key error: %w", err)
			}

			return &KVSet{config: config, ModuleId: moduleIdString, Key: keyString, logger: slog.Default().With("component", "processor", "type", config.Type)}, nil
		},
	})
}
