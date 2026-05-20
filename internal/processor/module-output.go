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

type ModuleOutput struct {
	config   config.ProcessorConfig
	ModuleId string
	logger   *slog.Logger
	module   common.OutputModule
}

func (ro *ModuleOutput) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {

	if ro.module == nil {
		if wrappedPayload.Modules == nil {
			wrappedPayload.End = true
			return wrappedPayload, errors.New("module.output wrapped payload has no modules")
		}

		module, ok := wrappedPayload.Modules[ro.ModuleId]
		if !ok {
			wrappedPayload.End = true
			return wrappedPayload, fmt.Errorf("module.output unable to find module with id: %s", ro.ModuleId)
		}

		outputModule, ok := module.(common.OutputModule)
		if !ok {
			wrappedPayload.End = true
			return wrappedPayload, fmt.Errorf("module.output module with id %s is not an OutputModule", ro.ModuleId)
		}
		ro.module = outputModule
	}

	err := ro.module.Output(ctx, wrappedPayload.Payload)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("module.output failed to send output: %w", err)
	}

	return wrappedPayload, nil
}

func (ro *ModuleOutput) Type() string {
	return ro.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "module.output",
		Title: "Module Output",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"module": {
					Title:       "Module ID",
					Type:        "string",
					Description: "ID of module to send output to",
				},
			},
			Required:             []string{"module"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ProcessorConfig) (Processor, error) {

			params := config.Params

			moduleId, err := params.GetString("module")

			if err != nil {
				return nil, fmt.Errorf("module.output module error: %w", err)
			}

			return &ModuleOutput{config: config, ModuleId: moduleId, logger: slog.Default().With("component", "processor", "type", config.Type)}, nil
		},
	})
}
