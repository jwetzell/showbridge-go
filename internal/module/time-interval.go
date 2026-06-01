package module

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "time.interval",
		Title: "Interval",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"duration": {
					Title:       "Duration",
					Description: "time in milliseconds between emitted events",
					Type:        "integer",
				},
			},
			Required:             []string{"duration"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ModuleConfig) (common.Module, error) {
			params := config.Params

			durationInt, err := params.GetInt("duration")
			if err != nil {
				return nil, fmt.Errorf("time.interval duration error: %w", err)
			}
			return &TimeInterval{Duration: uint32(durationInt), config: config, logger: CreateLogger(config)}, nil
		},
	})
}

type TimeInterval struct {
	config       config.ModuleConfig
	Duration     uint32
	ctx          context.Context
	inputHandler common.InputHandler
	ticker       *time.Ticker
	logger       *slog.Logger
	cancel       context.CancelFunc
}

func (i *TimeInterval) Id() string {
	return i.config.Id
}

func (i *TimeInterval) Type() string {
	return i.config.Type
}

func (i *TimeInterval) Start(ctx context.Context, inputHandler common.InputHandler) error {
	i.logger.Debug("running")
	i.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	i.ctx = moduleContext
	i.cancel = cancel

	ticker := time.NewTicker(time.Millisecond * time.Duration(i.Duration))
	i.ticker = ticker

	for i.ctx.Err() == nil {
		select {
		case <-ticker.C:
			if i.inputHandler != nil {
				i.inputHandler(i.ctx, i.Id(), time.Now())
			}
		default:
			continue
		}
	}
	<-i.ctx.Done()
	i.logger.Debug("done")
	return nil
}

func (i *TimeInterval) Stop() {
	if i.cancel != nil {
		defer i.cancel()
	}
	if i.ticker != nil {
		i.ticker.Stop()
	}
}
