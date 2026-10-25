package processor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "module.output",
		Title: "Module Output",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"module": {
					Title:       "Module ID",
					Description: "ID of module to send output to",
					Type:        "string",
				},
			},
			Required:             []string{"module"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ProcessorConfig) (routing.Processor, error) {

			processorParams := config.Params

			moduleId, err := processorParams.GetString("module")

			if err != nil {
				return nil, fmt.Errorf("module.output module error: %w", err)
			}

			return &ModuleOutput{config: config, ModuleId: moduleId, logger: slog.Default().With("component", "processor", "type", config.Type)}, nil
		},
	})
}

type ModuleOutput struct {
	config   config.ProcessorConfig
	ModuleId string
	logger   *slog.Logger
	module   routing.OutputModule
}

func (mo *ModuleOutput) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, inputHandler routing.InputHandler) (routing.WrappedPayload, error) {

	if mo.module == nil {
		if moduleManager == nil {
			wrappedPayload.End = true
			return wrappedPayload, errors.New("module.output module manager is nil")
		}

		module, ok := moduleManager.GetModule(mo.ModuleId)
		if !ok {
			wrappedPayload.End = true
			return wrappedPayload, fmt.Errorf("module.output unable to find module with id: %s", mo.ModuleId)
		}

		outputModule, ok := module.(routing.OutputModule)
		if !ok {
			wrappedPayload.End = true
			return wrappedPayload, fmt.Errorf("module.output module with id %s is not an OutputModule", mo.ModuleId)
		}
		mo.module = outputModule
	}

	err := mo.module.Output(ctx, wrappedPayload.Payload)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("module.output failed to send output: %w", err)
	}

	return wrappedPayload, nil
}
