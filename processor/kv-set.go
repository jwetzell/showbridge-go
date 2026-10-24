package processor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/common"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/payload"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "kv.set",
		Title: "Set Key",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"module": {
					Title:       "Module ID",
					Description: "ID of the key-value module to set the value in",
					Type:        "string",
				},
				"key": {
					Title:       "Key",
					Description: "key to set in the key-value module",
					Type:        "string",
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

type KVSet struct {
	config   config.ProcessorConfig
	ModuleId string
	Key      string
	logger   *slog.Logger
	module   common.KeyValueModule
}

func (kvs *KVSet) Process(ctx context.Context, wrappedPayload payload.WrappedPayload, moduleManager common.ModuleManager, inputHandler common.InputHandler) (payload.WrappedPayload, error) {
	if kvs.module == nil {
		if moduleManager == nil {
			wrappedPayload.End = true
			return wrappedPayload, errors.New("kv.set module manager is nil")
		}

		module, ok := moduleManager.GetModule(kvs.ModuleId)
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

	err := kvs.module.Set(ctx, kvs.Key, wrappedPayload.Payload)
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("kv.set error setting key: %w", err)
	}

	return wrappedPayload, nil
}
