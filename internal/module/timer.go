package module

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/route"
)

type Timer struct {
	config   config.ModuleConfig
	Duration uint32
	ctx      context.Context
	router   route.RouteIO
	timer    *time.Timer
	logger   *slog.Logger
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "gen.timer",
		New: func(ctx context.Context, config config.ModuleConfig, router route.RouteIO) (Module, error) {
			params := config.Params

			duration, ok := params["duration"]
			if !ok {
				return nil, errors.New("gen.timer requires a duration parameter")
			}

			durationNum, ok := duration.(float64)

			if !ok {
				return nil, errors.New("gen.timer duration must be a number")
			}

			return &Timer{Duration: uint32(durationNum), config: config, ctx: ctx, router: router, logger: CreateLogger(config)}, nil
		},
	})
}

func (t *Timer) Id() string {
	return t.config.Id
}

func (t *Timer) Type() string {
	return t.config.Type
}

func (t *Timer) Run() error {
	t.timer = time.NewTimer(time.Millisecond * time.Duration(t.Duration))
	defer t.timer.Stop()
	for {
		select {
		case <-t.ctx.Done():
			t.timer.Stop()
			t.logger.Debug("done")
			return nil
		case time := <-t.timer.C:
			if t.router != nil {
				t.router.HandleInput(t.Id(), time)
			}
		}
	}
}

func (t *Timer) Output(payload any) error {
	t.timer.Reset(time.Millisecond * time.Duration(t.Duration))
	return nil
}
