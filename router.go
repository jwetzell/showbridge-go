package showbridge

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
)

type RoutingError struct {
	Index int
	Error error
}

type Router struct {
	contextCancel   context.CancelFunc
	Context         context.Context
	ModuleInstances []Module
	RouteInstances  []*Route
	moduleWait      sync.WaitGroup
}

func NewRouter(ctx context.Context, config Config) (*Router, []ModuleError, []RouteError) {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	slog.Debug("creating router")

	routerContext, cancel := context.WithCancel(ctx)
	router := Router{
		Context:         routerContext,
		contextCancel:   cancel,
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
			continue
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
		moduleInstance.RegisterRouter(r)
		r.moduleWait.Add(1)
		go func() {
			err := moduleInstance.Run()
			if err != nil {
				slog.Error("error encountered running module", "id", moduleInstance.Id(), "error", err)
			}
			r.moduleWait.Done()
		}()
	}
	<-r.Context.Done()
	r.moduleWait.Wait()
	slog.Info("router context done")
}

func (r *Router) Stop() {
	r.contextCancel()
}

func (r *Router) HandleInput(sourceId string, payload any) []RoutingError {
	var routingErrors []RoutingError
	for routeIndex, route := range r.RouteInstances {
		if route.Input == sourceId {
			err := route.HandleInput(sourceId, payload)
			if err != nil {
				if routingErrors == nil {
					routingErrors = []RoutingError{}
				}
				routingErrors = append(routingErrors, RoutingError{
					Index: routeIndex,
					Error: err,
				})
				slog.Error("router unable to route input", "route", routeIndex, "source", sourceId, "error", err)
			}
		}
	}
	return routingErrors
}

func (r *Router) HandleOutput(destinationId string, payload any) error {
	for _, moduleInstance := range r.ModuleInstances {
		if moduleInstance.Id() == destinationId {
			return moduleInstance.Output(payload)
		}
	}
	return fmt.Errorf("no module instance found for destination %s", destinationId)
}
