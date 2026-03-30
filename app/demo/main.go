package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jwetzell/showbridge-go"
	"github.com/jwetzell/showbridge-go/internal/config"
)

func main() {

	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.SetDefault(slog.New(slog.NewTextHandler(NewLogWriter("logs"), &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	router, moduleConfigErrors, routeConfigErrors := showbridge.NewRouter(config.Config{
		Api: config.ApiConfig{
			Enabled: false,
			Port:    0,
		},
		Modules: []config.ModuleConfig{
			{
				Id:   "timer",
				Type: "time.interval",
				Params: map[string]any{
					"duration": 1000,
				},
			},
			{
				Id:   "button1",
				Type: "web.onclick",
				Params: map[string]any{
					"id": "button1",
				},
			},
			{
				Id:   "button2",
				Type: "web.onclick",
				Params: map[string]any{
					"id": "button2",
				},
			},
		},
		Routes: []config.RouteConfig{
			{
				Input: "timer",
				Processors: []config.ProcessorConfig{
					{
						Type: "debug.log",
					},
				},
			},
			{
				Input: "button1",
				Processors: []config.ProcessorConfig{
					{
						Type: "string.create",
						Params: map[string]any{
							"template": "{{.Payload.UnixMilli}}",
						},
					},
					{
						Type: "debug.log",
					},
					{
						Type: "web.set",
						Params: map[string]any{
							"id":       "output1",
							"property": "innerText",
							"value":    "Button1 Pressed @ {{.Payload}}",
						},
					},
				},
			},
			{
				Input: "button2",
				Processors: []config.ProcessorConfig{
					{
						Type: "string.create",
						Params: map[string]any{
							"template": "{{.Payload.UnixMilli}}",
						},
					},
					{
						Type: "debug.log",
					},
					{
						Type: "web.set",
						Params: map[string]any{
							"id":       "output2",
							"property": "innerText",
							"value":    "Button2 Pressed @ {{.Payload}}",
						},
					},
				},
			},
		},
	})

	if len(moduleConfigErrors) > 0 {
		for _, err := range moduleConfigErrors {
			println("Module config error:", err.Error)
		}
	}

	if len(routeConfigErrors) > 0 {
		for _, err := range routeConfigErrors {
			println("Route config error:", err.Error)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		router.Start(ctx)
		fmt.Println("router stopped")
	}()
	<-ctx.Done()
}
