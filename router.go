package showbridge

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type Router struct {
	Context         context.Context
	ModuleInstances []Module
}

func NewRouter(ctx context.Context, config Config) (*Router, error) {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	slog.Debug("creating router", "config", config)

	router := Router{
		Context:         ctx,
		ModuleInstances: []Module{},
	}

	for _, moduleDecl := range config.Modules {

		moduleInfo, ok := moduleRegistry[moduleDecl.Type]
		if !ok {
			return nil, fmt.Errorf("problem loading module registration for module type: %s", moduleDecl.Type)
		}

		moduleInstance, err := moduleInfo.New(moduleDecl.Params)
		if err != nil {
			return nil, err
		}

		router.ModuleInstances = append(router.ModuleInstances, moduleInstance)

	}

	return &router, nil
}

func (r *Router) Run() {
	for _, moduleInstance := range r.ModuleInstances {
		go moduleInstance.Run(r.Context)
	}
	<-r.Context.Done()
}
