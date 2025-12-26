package module

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/route"
)

type TimeTimer struct {
	config   config.ModuleConfig
	Duration uint32
	ctx      context.Context
	router   route.RouteIO
	timer    *time.Timer
	logger   *slog.Logger
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "time.timer",
		New: func(ctx context.Context, config config.ModuleConfig, router route.RouteIO) (Module, error) {
			params := config.Params

			duration, ok := params["duration"]
			if !ok {
				return nil, errors.New("time.timer requires a duration parameter")
			}

			durationNum, ok := duration.(float64)

			if !ok {
				return nil, errors.New("time.timer duration must be a number")
			}

			return &TimeTimer{Duration: uint32(durationNum), config: config, ctx: ctx, router: router, logger: CreateLogger(config)}, nil
		},
	})
}

func (t *TimeTimer) Id() string {
	return t.config.Id
}

func (t *TimeTimer) Type() string {
	return t.config.Type
}

func (t *TimeTimer) Run() error {
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

func (t *TimeTimer) Output(payload any) error {
	t.timer.Reset(time.Millisecond * time.Duration(t.Duration))
	return nil
}
