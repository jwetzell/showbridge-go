package processor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type RouterOutput struct {
	config   config.ProcessorConfig
	ModuleId string
	logger   *slog.Logger
}

func (dl *RouterOutput) Process(ctx context.Context, payload any) (any, error) {

	router, ok := ctx.Value(common.RouterContextKey).(common.RouteIO)
	if !ok {
		return nil, errors.New("router.output no router found")
	}

	err := router.HandleOutput(ctx, dl.ModuleId, payload)

	if err != nil {
		return nil, fmt.Errorf("router.output failed to send output: %w", err)
	}

	return payload, nil
}

func (dl *RouterOutput) Type() string {
	return dl.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "router.output",
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
