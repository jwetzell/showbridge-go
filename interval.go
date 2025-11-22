package showbridge

import (
	"fmt"
	"log/slog"
	"time"
)

type Interval struct {
	config   ModuleConfig
	Duration uint32
	router   *Router
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "gen.interval",
		New: func(config ModuleConfig) (Module, error) {
			params := config.Params

			duration, ok := params["duration"]
			if !ok {
				return nil, fmt.Errorf("interval requires a duration parameter")
			}

			durationNum, ok := duration.(float64)

			if !ok {
				return nil, fmt.Errorf("interval duration must be number")
			}

			return &Interval{Duration: uint32(durationNum), config: config}, nil
		},
	})
}

func (i *Interval) Id() string {
	return i.config.Id
}

func (i *Interval) Type() string {
	return i.config.Type
}

func (i *Interval) RegisterRouter(router *Router) {
	i.router = router
}

func (i *Interval) Run() error {
	ticker := time.NewTicker(time.Millisecond * time.Duration(i.Duration))
	defer ticker.Stop()
	for {
		select {
		case <-i.router.Context.Done():
			ticker.Stop()
			slog.Debug("router context done in module", "id", i.config.Id)
			return nil
		case t := <-ticker.C:
			if i.router != nil {
				i.router.HandleInput(i.config.Id, t)
			}
		}
	}

}

func (i *Interval) Output(payload any) error {
	return fmt.Errorf("interval output is not implemented")
}
