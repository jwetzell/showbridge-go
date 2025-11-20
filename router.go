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
	RouteInstances  []*Route
}

func NewRouter(ctx context.Context, config Config) (*Router, error) {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(logger)

	slog.Debug("creating router")

	router := Router{
		Context:         ctx,
		ModuleInstances: []Module{},
		RouteInstances:  []*Route{},
	}

	for _, moduleDecl := range config.Modules {

		moduleInfo, ok := moduleRegistry[moduleDecl.Type]
		if !ok {
			return nil, fmt.Errorf("problem loading module registration for module type: %s", moduleDecl.Type)
		}

		moduleInstanceExists := false
		for _, moduleInstance := range router.ModuleInstances {
			if moduleInstance.Id() == moduleDecl.Id {
				moduleInstanceExists = true
				slog.Warn("module id conflict", "id", moduleDecl.Id, "type", moduleDecl.Type)
				break
			}
		}

		if !moduleInstanceExists {
			moduleInstance, err := moduleInfo.New(moduleDecl)
			if err != nil {
				return nil, err
			}

			router.ModuleInstances = append(router.ModuleInstances, moduleInstance)
		}

	}

	for routeIndex, routeDecl := range config.Routes {
		router.RouteInstances = append(router.RouteInstances, NewRoute(routeIndex, routeDecl, &router))
	}

	for _, moduleInstance := range router.ModuleInstances {
		moduleInstance.RegisterRouter(&router)
	}

	return &router, nil
}

func (r *Router) Run() {
	for _, moduleInstance := range r.ModuleInstances {
		go moduleInstance.Run(r.Context)
	}
	<-r.Context.Done()
}

func (r *Router) HandleInput(sourceId string, payload any) {
	for _, route := range r.RouteInstances {
		if route.Input == sourceId {
			route.HandleInput(sourceId, payload)
		}
	}
}

func (r *Router) HandleOutput(destinationId string, payload any) error {
	for _, moduleInstance := range r.ModuleInstances {
		if moduleInstance.Id() == destinationId {
			return moduleInstance.Output(payload)
		}
	}
	return fmt.Errorf("no module instance found for destination %s", destinationId)
}
