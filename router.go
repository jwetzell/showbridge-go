package showbridge

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
	"github.com/jwetzell/showbridge-go/internal/route"
)

type Router struct {
	contextCancel   context.CancelFunc
	Context         context.Context
	ModuleInstances []module.Module
	RouteInstances  []route.Route
	moduleWait      sync.WaitGroup
}

func NewRouter(ctx context.Context, config config.Config) (*Router, []module.ModuleError, []route.RouteError) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(logger)

	slog.Debug("creating router")

	routerContext, cancel := context.WithCancel(ctx)
	router := Router{
		Context:         routerContext,
		contextCancel:   cancel,
		ModuleInstances: []module.Module{},
		RouteInstances:  []route.Route{},
	}

	var moduleErrors []module.ModuleError

	for moduleIndex, moduleDecl := range config.Modules {

		moduleInfo, ok := module.ModuleRegistry[moduleDecl.Type]
		if !ok {
			if moduleErrors == nil {
				moduleErrors = []module.ModuleError{}
			}
			moduleErrors = append(moduleErrors, module.ModuleError{
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
					moduleErrors = []module.ModuleError{}
				}
				moduleErrors = append(moduleErrors, module.ModuleError{
					Index:  moduleIndex,
					Config: moduleDecl,
					Error:  fmt.Errorf("duplicate module id"),
				})
				break
			}
		}

		if !moduleInstanceExists {
			moduleInstance, err := moduleInfo.New(router.Context, moduleDecl, &router)
			if err != nil {
				if moduleErrors == nil {
					moduleErrors = []module.ModuleError{}
				}
				moduleErrors = append(moduleErrors, module.ModuleError{
					Index:  moduleIndex,
					Config: moduleDecl,
					Error:  err,
				})
				continue
			}

			router.ModuleInstances = append(router.ModuleInstances, moduleInstance)
		}

	}

	var routeErrors []route.RouteError
	for routeIndex, routeDecl := range config.Routes {
		routeInstance, err := route.NewRoute(routeDecl)
		if err != nil {
			if routeErrors == nil {
				routeErrors = []route.RouteError{}
			}
			routeErrors = append(routeErrors, route.RouteError{
				Index:  routeIndex,
				Config: routeDecl,
				Error:  err,
			})
			continue
		}
		router.RouteInstances = append(router.RouteInstances, routeInstance)
	}

	return &router, moduleErrors, routeErrors
}

func (r *Router) Run() {
	slog.Info("running router")
	for _, moduleInstance := range r.ModuleInstances {
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
	slog.Info("router done")
}

func (r *Router) Stop() {
	r.contextCancel()
}

func (r *Router) HandleInput(sourceId string, payload any) []route.RouteIOError {
	var routingErrors []route.RouteIOError
	for routeIndex, routeInstance := range r.RouteInstances {
		if routeInstance.Input() == sourceId {
			err := routeInstance.HandleInput(r.Context, sourceId, payload, r)
			if err != nil {
				if routingErrors == nil {
					routingErrors = []route.RouteIOError{}
				}
				routingErrors = append(routingErrors, route.RouteIOError{
					Index: routeIndex,
					Error: err,
				})
				slog.Error("router unable to route input", "route", routeIndex, "source", sourceId, "error", err)
			}
		}
	}
	return routingErrors
}

func (r *Router) HandleOutput(sourceId string, destinationId string, payload any) error {
	for _, moduleInstance := range r.ModuleInstances {
		if moduleInstance.Id() == destinationId {
			return moduleInstance.Output(payload)
		}
	}
	return fmt.Errorf("router could not find module instance for destination %s", destinationId)
}
