package showbridge

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type Timer struct {
	config   config.ModuleConfig
	Duration uint32
	router   *Router
	timer    *time.Timer
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "gen.timer",
		New: func(config config.ModuleConfig, router *Router) (Module, error) {
			params := config.Params

			duration, ok := params["duration"]
			if !ok {
				return nil, fmt.Errorf("gen.timer requires a duration parameter")
			}

			durationNum, ok := duration.(float64)

			if !ok {
				return nil, fmt.Errorf("gen.timer duration must be a number")
			}

			return &Timer{Duration: uint32(durationNum), config: config, router: router}, nil
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
		case <-t.router.Context.Done():
			t.timer.Stop()
			slog.Debug("router context done in module", "id", t.config.Id)
			return nil
		case time := <-t.timer.C:
			if t.router != nil {
				t.router.HandleInput(t.config.Id, time)
			}
		}
	}
}

func (t *Timer) Output(payload any) error {
	t.timer.Reset(time.Millisecond * time.Duration(t.Duration))
	return nil
}
