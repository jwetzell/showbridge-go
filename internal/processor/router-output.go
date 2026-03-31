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

type RouterOutput struct {
	config   config.ProcessorConfig
	ModuleId string
	logger   *slog.Logger
}

func (ro *RouterOutput) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {

	router, ok := ctx.Value(common.RouterContextKey).(common.RouteIO)
	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("router.output no router found")
	}

	err := router.HandleOutput(ctx, ro.ModuleId, wrappedPayload.Payload)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("router.output failed to send output: %w", err)
	}

	return wrappedPayload, nil
}

func (ro *RouterOutput) Type() string {
	return ro.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "router.output",
		Title: "Router Output",
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
				return nil, fmt.Errorf("router.output module error: %w", err)
			}

			return &RouterOutput{config: config, ModuleId: moduleId, logger: slog.Default().With("component", "processor", "type", config.Type)}, nil
		},
	})
}
