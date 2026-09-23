package module

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/robfig/cron/v3"
)

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "time.cron",
		Title: "Cron",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"cron": {
					Title:       "Cron Spec",
					Description: "cron expression specifying the schedule",
					Type:        "string",
				},
			},
			Required:             []string{"cron"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ModuleConfig) (common.Module, error) {
			params := config.Params

			cronString, err := params.GetString("cron")
			if err != nil {
				return nil, fmt.Errorf("time.cron cron error: %w", err)
			}

			return &TimeCron{CronString: cronString, config: config, logger: CreateLogger(config)}, nil
		},
	})
}

type TimeCron struct {
	config       config.ModuleConfig
	CronString   string
	ctx          context.Context
	inputHandler common.InputHandler
	logger       *slog.Logger
	cancel       context.CancelFunc
	cron         *cron.Cron
}

func (t *TimeCron) Id() string {
	return t.config.Id
}

func (t *TimeCron) Type() string {
	return t.config.Type
}

func (t *TimeCron) Start(ctx context.Context, inputHandler common.InputHandler) error {
	t.logger.Debug("running")
	t.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	t.ctx = moduleContext
	t.cancel = cancel

	t.cron = cron.New()

	t.cron.AddFunc(t.CronString, func() {
		if t.inputHandler != nil {
			t.inputHandler(t.ctx, t.Id(), time.Now())
		}
	})
	t.cron.Start()

	<-t.ctx.Done()
	t.logger.Debug("done")
	return nil
}

func (t *TimeCron) Stop() {
	if t.cancel != nil {
		defer t.cancel()
	}
	if t.cron != nil {
		t.cron.Stop()
	}
}
