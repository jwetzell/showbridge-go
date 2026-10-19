//go:build cgo

package module

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/0xcafed00d/joystick"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/common"
)

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "joystick.input",
		Title: "Joystick Input",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"id": {
					Title:       "Id",
					Description: "the name of the MIDI port to listen to",
					Type:        "number",
				},
			},
			Required:             []string{"id"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ModuleConfig) (common.Module, error) {
			params := config.Params
			id, err := params.GetInt("id")
			if err != nil {
				return nil, fmt.Errorf("joystick.input id error: %w", err)
			}

			return &JoystickInput{config: config, JoystickId: id, logger: CreateLogger(config)}, nil
		},
	})
}

type JoystickInput struct {
	config       config.ModuleConfig
	ctx          context.Context
	inputHandler common.InputHandler
	JoystickId   int
	logger       *slog.Logger
	cancel       context.CancelFunc
	js           joystick.Joystick
}

func (mi *JoystickInput) Id() string {
	return mi.config.Id
}

func (mi *JoystickInput) Type() string {
	return mi.config.Type
}

func (mi *JoystickInput) Start(ctx context.Context, inputHandler common.InputHandler) error {
	mi.logger.Debug("running")
	mi.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	mi.ctx = moduleContext
	mi.cancel = cancel

	js, err := joystick.Open(mi.JoystickId)
	if err != nil {
		return fmt.Errorf("joystick.input can't open input port: %d", mi.JoystickId)
	}
	mi.js = js

	go func() {
		for {
			select {
			case <-mi.ctx.Done():
				return
			default:
				event, err := mi.js.Read()
				if err != nil {
					mi.logger.Error("error reading joystick input", "error", err)
					continue
				}
				if mi.inputHandler != nil {
					mi.inputHandler(mi.ctx, mi.Id(), event)
				}
			}
		}
	}()

	// if err != nil {
	// 	return err
	// }
	// mi.stop = stop

	<-mi.ctx.Done()
	mi.logger.Debug("done")
	return nil
}

func (mi *JoystickInput) Stop() {
	if mi.cancel != nil {
		defer mi.cancel()
	}
	mi.cancel()
	mi.js.Close()
}
