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

type MIDIInput struct {
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
		Type: "midi.input",
		New: func(config config.ModuleConfig) (common.Module, error) {
			params := config.Params
			portString, err := params.GetString("port")
			if err != nil {
				return nil, fmt.Errorf("midi.input port error: %w", err)
			}

			return &MIDIInput{config: config, Port: portString, logger: CreateLogger(config)}, nil
		},
	})
}

func (mi *MIDIInput) Id() string {
	return mi.config.Id
}

func (mi *MIDIInput) Type() string {
	return mi.config.Type
}

func (mi *MIDIInput) Start(ctx context.Context) error {
	mi.logger.Debug("running")
	defer midi.CloseDriver()
	router, ok := ctx.Value(common.RouterContextKey).(common.RouteIO)

	if !ok {
		return errors.New("midi.input unable to get router from context")
	}
	mi.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	mi.ctx = moduleContext
	mi.cancel = cancel

	in, err := midi.FindInPort(mi.Port)
	if err != nil {
		return fmt.Errorf("midi.input can't find input port: %s", mi.Port)
	}

	stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
		if mi.router != nil {
			mi.router.HandleInput(mi.ctx, mi.Id(), msg)
		}
	}, midi.UseSysEx())

	if err != nil {
		return err
	}

	defer stop()

	<-mi.ctx.Done()
	mi.logger.Debug("done")
	return nil
}

func (mi *MIDIInput) Output(ctx context.Context, payload any) error {
	return errors.New("midi.input output is not implemented")
}

func (mi *MIDIInput) Stop() {
	mi.cancel()
}

func (mi *MIDIInput) Get(key string) (any, error) {
	switch key {
	case "port":
		return mi.Port, nil
	default:
		return nil, errors.New("midi.input key not found")
	}
}
