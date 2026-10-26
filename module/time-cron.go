package module

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/routing"
	"github.com/robfig/cron/v3"
)

func init() {
	RegisterModule(ModuleRegistration{
		Type:        "time.cron",
		Title:       "Cron",
		Description: "Triggers events based on the specified cron schedule",
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
		New: func(config config.ModuleConfig) (routing.Module, error) {
			moduleParams := config.Params

			cronString, err := moduleParams.GetString("cron")
			if err != nil {
				return nil, fmt.Errorf("time.cron cron error: %w", err)
			}

			return &TimeCron{CronString: cronString, config: config, logger: CreateLogger(config)}, nil
		},
	})
}

type TimeCron struct {
	config     config.ModuleConfig
	CronString string
	ctx        context.Context
	router     routing.Router
	logger     *slog.Logger
	cancel     context.CancelFunc
	cron       *cron.Cron
}

func (t *TimeCron) Start(ctx context.Context, router routing.Router) error {
	t.logger.Debug("running")
	t.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	t.ctx = moduleContext
	t.cancel = cancel

	t.cron = cron.New()

	t.cron.AddFunc(t.CronString, func() {
		if t.router != nil {
			t.router.HandleInput(t.ctx, t.config.Id, time.Now())
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
