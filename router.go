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

func NewRouter(ctx context.Context, config Config) (*Router, []ModuleError, []RouteError) {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	slog.Debug("creating router")

	router := Router{
		Context:         ctx,
		ModuleInstances: []Module{},
		RouteInstances:  []*Route{},
	}

	var moduleErrors []ModuleError

	for moduleIndex, moduleDecl := range config.Modules {

		moduleInfo, ok := moduleRegistry[moduleDecl.Type]
		if !ok {
			if moduleErrors == nil {
				moduleErrors = []ModuleError{}
			}
			moduleErrors = append(moduleErrors, ModuleError{
				Index:  moduleIndex,
				Config: moduleDecl,
				Error:  fmt.Errorf("module type not defined"),
			})

		}

		moduleInstanceExists := false
		for _, moduleInstance := range router.ModuleInstances {
			if moduleInstance.Id() == moduleDecl.Id {
				moduleInstanceExists = true
				if moduleErrors == nil {
					moduleErrors = []ModuleError{}
				}
				moduleErrors = append(moduleErrors, ModuleError{
					Index:  moduleIndex,
					Config: moduleDecl,
					Error:  fmt.Errorf("duplicate module id"),
				})
				break
			}
		}

		if !moduleInstanceExists {
			moduleInstance, err := moduleInfo.New(moduleDecl)
			if err != nil {
				if moduleErrors == nil {
					moduleErrors = []ModuleError{}
				}
				moduleErrors = append(moduleErrors, ModuleError{
					Index:  moduleIndex,
					Config: moduleDecl,
					Error:  err,
				})
				continue
			}

			router.ModuleInstances = append(router.ModuleInstances, moduleInstance)
		}

	}

	var routeErrors []RouteError
	for routeIndex, routeDecl := range config.Routes {
		route, err := NewRoute(routeIndex, routeDecl, &router)
		if err != nil {
			if routeErrors == nil {
				routeErrors = []RouteError{}
			}
			routeErrors = append(routeErrors, RouteError{
				Index:  routeIndex,
				Config: routeDecl,
				Error:  err,
			})
			continue
		}
		router.RouteInstances = append(router.RouteInstances, route)
	}

	for _, moduleInstance := range router.ModuleInstances {
		slog.Debug("registering router with module", "id", moduleInstance.Id())
		moduleInstance.RegisterRouter(&router)
	}

	return &router, moduleErrors, routeErrors
}

func (r *Router) Run() {
	for _, moduleInstance := range r.ModuleInstances {
		go moduleInstance.Run(r.Context)
	}
	<-r.Context.Done()
}

func (r *Router) HandleInput(sourceId string, payload any) {
	for routeIndex, route := range r.RouteInstances {
		if route.Input == sourceId {
			err := route.HandleInput(sourceId, payload)
			if err != nil {
				slog.Error("router unable to route input", "route", routeIndex, "source", sourceId, "error", err)
			}
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
