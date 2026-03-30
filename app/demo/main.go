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
				Id:   "midi",
				Type: "midi.input",
				Params: map[string]any{
					"port": "Launchpad S",
				},
			},
		},
		Routes: []config.RouteConfig{
			{
				Input: "midi",
				Processors: []config.ProcessorConfig{
					{
						Type: "debug.log",
					},
					{
						Type: "web.set",
						Params: map[string]any{
							"id":       "output",
							"property": "textContent",
							"value":    "{{.Payload}}",
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
