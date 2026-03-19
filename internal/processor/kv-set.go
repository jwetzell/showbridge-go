package processor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"log/slog"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type KVSet struct {
	config   config.ProcessorConfig
	ModuleId string
	Key      string
	Value    *template.Template
	logger   *slog.Logger
}

func (kvs *KVSet) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {

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

	var valueBuffer bytes.Buffer
	err := kvs.Value.Execute(&valueBuffer, wrappedPayload)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	err = kvModule.Set(kvs.Key, valueBuffer.String())
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
		Type: "kv.set",
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

			valueString, err := params.GetString("value")
			if err != nil {
				return nil, fmt.Errorf("kv.set value error: %w", err)
			}
			valueTemplate, err := template.New("template").Parse(valueString)

			if err != nil {
				return nil, err
			}

			return &KVSet{config: config, ModuleId: moduleIdString, Key: keyString, Value: valueTemplate, logger: slog.Default().With("component", "processor", "type", config.Type)}, nil
		},
	})
}
