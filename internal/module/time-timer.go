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
		Type:  "time.timer",
		Title: "Timer",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"duration": {
					Title:       "Duration",
					Type:        "integer",
					Description: "Interval duration in milliseconds",
				},
			},
			Required:             []string{"duration"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ModuleConfig) (common.Module, error) {
			params := config.Params

			durationNum, err := params.GetInt("duration")
			if err != nil {
				return nil, fmt.Errorf("time.timer duration error: %w", err)
			}

			return &TimeTimer{Duration: uint32(durationNum), config: config, logger: CreateLogger(config)}, nil
		},
	})
}

type TimeTimer struct {
	config       config.ModuleConfig
	Duration     uint32
	ctx          context.Context
	inputHandler common.InputHandler
	timer        *time.Timer
	logger       *slog.Logger
	cancel       context.CancelFunc
}

func (t *TimeTimer) Id() string {
	return t.config.Id
}

func (t *TimeTimer) Type() string {
	return t.config.Type
}

func (t *TimeTimer) Start(ctx context.Context, inputHandler common.InputHandler) error {
	t.logger.Debug("running")
	t.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	t.ctx = moduleContext
	t.cancel = cancel

	t.timer = time.NewTimer(time.Millisecond * time.Duration(t.Duration))
	for {
		select {
		case <-t.ctx.Done():
			return nil
		case time := <-t.timer.C:
			if t.inputHandler != nil {
				t.inputHandler(t.ctx, t.Id(), time)
			}
		}
	}
}

func (t *TimeTimer) Stop() {
	if t.cancel != nil {
		t.cancel()
	}
	if t.timer != nil {
		t.timer.Stop()
		t.timer = nil
	}
	t.logger.Debug("done")
}
