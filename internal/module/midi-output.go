//go:build cgo

package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

type MIDIOutput struct {
	config   config.ModuleConfig
	ctx      context.Context
	router   common.RouteIO
	Port     string
	SendFunc func(midi.Message) error
	logger   *slog.Logger
	cancel   context.CancelFunc
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "midi.output",
		New: func(config config.ModuleConfig) (common.Module, error) {
			params := config.Params

			portString, err := params.GetString("port")
			if err != nil {
				return nil, fmt.Errorf("midi.output port error: %w", err)
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

func (mo *MIDIOutput) Start(ctx context.Context) error {
	mo.logger.Debug("running")
	defer midi.CloseDriver()
	router, ok := ctx.Value(common.RouterContextKey).(common.RouteIO)

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

	payloadMessage, ok := common.GetAnyAs[midi.Message](payload)

	if !ok {
		return errors.New("midi.output can only ouptut midi.Message")
	}

	return mo.SendFunc(payloadMessage)
}

func (mo *MIDIOutput) Stop() {
	mo.cancel()
}

func (mo *MIDIOutput) Get(key string) (any, error) {
	switch key {
	case "port":
		return mo.Port, nil
	default:
		return nil, errors.New("midi.output key not found")
	}
}
