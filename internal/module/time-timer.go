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
	cancel   context.CancelFunc
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "time.timer",
		New: func(config config.ModuleConfig) (Module, error) {
			params := config.Params

			duration, ok := params["duration"]
			if !ok {
				return nil, errors.New("time.timer requires a duration parameter")
			}

			durationNum, ok := duration.(float64)

			if !ok {
				return nil, errors.New("time.timer duration must be a number")
			}

			return &TimeTimer{Duration: uint32(durationNum), config: config, logger: CreateLogger(config)}, nil
		},
	})
}

func (t *TimeTimer) Id() string {
	return t.config.Id
}

func (t *TimeTimer) Type() string {
	return t.config.Type
}

func (t *TimeTimer) Start(ctx context.Context) error {
	t.logger.Debug("running")
	router, ok := ctx.Value(route.RouterContextKey).(route.RouteIO)

	if !ok {
		return errors.New("net.tcp.client unable to get router from context")
	}
	t.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	t.ctx = moduleContext
	t.cancel = cancel

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
				t.router.HandleInput(t.ctx, t.Id(), time)
			}
		}
	}
}

func (t *TimeTimer) Output(ctx context.Context, payload any) error {
	t.timer.Reset(time.Millisecond * time.Duration(t.Duration))
	return nil
}

func (t *TimeTimer) Stop() {
	t.cancel()
}
