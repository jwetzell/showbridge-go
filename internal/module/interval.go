package module

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/route"
)

type Interval struct {
	config   config.ModuleConfig
	Duration uint32
	ctx      context.Context
	router   route.RouteIO
	ticker   *time.Ticker
	logger   *slog.Logger
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "gen.interval",
		New: func(ctx context.Context, config config.ModuleConfig, router route.RouteIO) (Module, error) {
			params := config.Params

			duration, ok := params["duration"]
			if !ok {
				return nil, fmt.Errorf("gen.interval requires a duration parameter")
			}

			durationNum, ok := duration.(float64)

			if !ok {
				return nil, fmt.Errorf("gen.interval duration must be number")
			}

			return &Interval{Duration: uint32(durationNum), config: config, ctx: ctx, router: router, logger: slog.Default().With("component", "module", "id", config.Id)}, nil
		},
	})
}

func (i *Interval) Id() string {
	return i.config.Id
}

func (i *Interval) Type() string {
	return i.config.Type
}

func (i *Interval) Run() error {
	ticker := time.NewTicker(time.Millisecond * time.Duration(i.Duration))
	i.ticker = ticker
	defer ticker.Stop()

	for {
		select {
		case <-i.ctx.Done():
			i.logger.Debug("router context done in module")
			return nil
		case <-ticker.C:
			if i.router != nil {
				i.router.HandleInput(i.Id(), time.Now())
			}
		}
	}

}

func (i *Interval) Output(payload any) error {
	i.ticker.Reset(time.Millisecond * time.Duration(i.Duration))
	return nil
}
