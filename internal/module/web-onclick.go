//go:build js

package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"syscall/js"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type WebOnClick struct {
	config    config.ModuleConfig
	ctx       context.Context
	router    common.RouteIO
	logger    *slog.Logger
	ElementId string
}

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "web.onclick",
		Title: "On Click",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"id": {
					Title:       "Element ID",
					Type:        "string",
					Description: "ID of the HTML element to attach the click listener to",
				},
			},
			Required:             []string{"duration"},
			AdditionalProperties: nil,
		},
		New: func(config config.ModuleConfig) (common.Module, error) {
			params := config.Params

			idString, err := params.GetString("id")
			if err != nil {
				return nil, fmt.Errorf("web.onclick id error: %w", err)
			}

			return &WebOnClick{ElementId: idString, config: config, logger: CreateLogger(config)}, nil
		},
	})
}

func (woc *WebOnClick) Id() string {
	return woc.config.Id
}

func (woc *WebOnClick) Type() string {
	return woc.config.Type
}

func (woc *WebOnClick) Start(ctx context.Context) error {
	woc.logger.Debug("running")
	router, ok := ctx.Value(common.RouterContextKey).(common.RouteIO)

	if !ok {
		return errors.New("net.tcp.client unable to get router from context")
	}
	woc.router = router
	woc.ctx = ctx

	element := js.Global().Get("document").Call("getElementById", woc.ElementId)

	if element.IsNull() || element.IsUndefined() {
		return fmt.Errorf("web.onclick unable to find element with id: %s", woc.ElementId)
	}

	element.Set("onclick", js.FuncOf(func(js.Value, []js.Value) interface{} {
		if woc.router != nil {
			woc.router.HandleInput(woc.ctx, woc.Id(), time.Now())
		}
		return nil
	}))

	<-ctx.Done()
	return nil
}

func (woc *WebOnClick) Stop() {
}
