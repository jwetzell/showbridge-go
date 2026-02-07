//go:build cgo

package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/route"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

type MIDIOutput struct {
	config   config.ModuleConfig
	ctx      context.Context
	router   route.RouteIO
	Port     string
	SendFunc func(midi.Message) error
	logger   *slog.Logger
	cancel   context.CancelFunc
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "midi.output",
		New: func(config config.ModuleConfig) (Module, error) {
			params := config.Params

			port, ok := params["port"]

			if !ok {
				return nil, errors.New("midi.output requires a port parameter")
			}

			portString, ok := port.(string)

			if !ok {
				return nil, errors.New("midi.output port must be a string")
			}

			return &MIDIOutput{config: config, Port: portString, logger: CreateLogger(config)}, nil
		},
	})
}

func (mo *MIDIOutput) Id() string {
	return mo.config.Id
}

func (mo *MIDIOutput) Type() string {
	return mo.config.Type
}

func (mo *MIDIOutput) Run(ctx context.Context) error {
	defer midi.CloseDriver()
	router, ok := ctx.Value(route.RouterContextKey).(route.RouteIO)

	if !ok {
		return errors.New("midi.output unable to get router from context")
	}
	mo.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	mo.ctx = moduleContext
	mo.cancel = cancel

	out, err := midi.FindOutPort(mo.Port)

	if err != nil {
		return fmt.Errorf("midi.output can't find output port: %s", mo.Port)
	}

	send, err := midi.SendTo(out)
	if err != nil {
		return err
	}

	mo.SendFunc = send

	<-mo.ctx.Done()
	mo.logger.Debug("done")
	return nil
}

func (mo *MIDIOutput) Output(ctx context.Context, payload any) error {
	if mo.SendFunc == nil {
		return errors.New("midi.output output is not setup")
	}

	payloadMessage, ok := payload.(midi.Message)

	if !ok {
		return errors.New("midi.output can only ouptut midi.Message")
	}

	return mo.SendFunc(payloadMessage)
}

func (mo *MIDIOutput) Stop() {
	mo.cancel()
}
