package showbridge

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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
	logger          *slog.Logger
}

func NewRouter(ctx context.Context, config config.Config) (*Router, []module.ModuleError, []route.RouteError) {

	routerContext, cancel := context.WithCancel(ctx)

	router := Router{
		contextCancel:   cancel,
		ModuleInstances: []module.Module{},
		RouteInstances:  []route.Route{},
		logger:          slog.Default().With("component", "router"),
	}

	router.Context = context.WithValue(routerContext, route.RouterContextKey, &router)

	router.logger.Debug("creating")

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
				Error:  errors.New("module type not defined"),
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
					Error:  errors.New("duplicate module id"),
				})
				break
			}
		}

		if !moduleInstanceExists {
			moduleInstance, err := moduleInfo.New(router.Context, moduleDecl)
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
	r.logger.Info("running")
	for _, moduleInstance := range r.ModuleInstances {
		r.moduleWait.Go(func() {
			err := moduleInstance.Run()
			if err != nil {
				r.logger.Error("error encountered running module", "error", err)
			}
		})
	}
	<-r.Context.Done()
	r.logger.Debug("waiting for modules to exit")
	r.moduleWait.Wait()
	r.logger.Info("done")
}

func (r *Router) Stop() {
	r.logger.Debug("stopping")
	r.contextCancel()
}

func (r *Router) HandleInput(sourceId string, payload any) []route.RouteIOError {
	var routingErrors []route.RouteIOError
	for routeIndex, routeInstance := range r.RouteInstances {
		if routeInstance.Input() == sourceId {
			err := routeInstance.HandleInput(context.WithValue(r.Context, route.SourceContextKey, sourceId), payload)
			if err != nil {
				if routingErrors == nil {
					routingErrors = []route.RouteIOError{}
				}
				routingErrors = append(routingErrors, route.RouteIOError{
					Index: routeIndex,
					Error: err,
				})
				r.logger.Error("unable to route input", "route", routeIndex, "source", sourceId, "error", err)
			}
		}
	}
	return routingErrors
}

func (r *Router) HandleOutput(ctx context.Context, destinationId string, payload any) error {
	for _, moduleInstance := range r.ModuleInstances {
		if moduleInstance.Id() == destinationId {
			return moduleInstance.Output(ctx, payload)
		}
	}
	return fmt.Errorf("router could not find module instance for destination %s", destinationId)
}
