package processor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

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
	ctxModules := ctx.Value(common.ModulesContextKey)
	if ctxModules == nil {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("kv.get unable to get modules from context")
	}

	moduleMap, ok := ctxModules.(map[string]common.Module)
	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("kv.get modules from context has wrong type")
	}

	module, ok := moduleMap[kvg.ModuleId]
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
		Type: "kv.get",
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
