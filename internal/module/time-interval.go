package module

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/route"
)

type TimeInterval struct {
	config   config.ModuleConfig
	Duration uint32
	ctx      context.Context
	router   route.RouteIO
	ticker   *time.Ticker
	logger   *slog.Logger
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "time.interval",
		New: func(ctx context.Context, config config.ModuleConfig, router route.RouteIO) (Module, error) {
			params := config.Params

			duration, ok := params["duration"]
			if !ok {
				return nil, errors.New("time.interval requires a duration parameter")
			}

			durationNum, ok := duration.(float64)

			if !ok {
				return nil, errors.New("time.interval duration must be number")
			}

			return &TimeInterval{Duration: uint32(durationNum), config: config, ctx: ctx, router: router, logger: CreateLogger(config)}, nil
		},
	})
}

func (i *TimeInterval) Id() string {
	return i.config.Id
}

func (i *TimeInterval) Type() string {
	return i.config.Type
}

func (i *TimeInterval) Run() error {
	ticker := time.NewTicker(time.Millisecond * time.Duration(i.Duration))
	i.ticker = ticker
	defer ticker.Stop()

	for {
		select {
		case <-i.ctx.Done():
			i.logger.Debug("done")
			return nil
		case <-ticker.C:
			if i.router != nil {
				i.router.HandleInput(i.Id(), time.Now())
			}
		}
	}

}

func (i *TimeInterval) Output(payload any) error {
	i.ticker.Reset(time.Millisecond * time.Duration(i.Duration))
	return nil
}
