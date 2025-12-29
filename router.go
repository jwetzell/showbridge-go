package showbridge

import (
	"context"
	"errors"
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

func (r *Router) HandleInput(ctx context.Context, sourceId string, payload any) (bool, []route.RouteIOError) {
	var routeIOErrors []route.RouteIOError
	routeFound := false

	var routeWaitGroup sync.WaitGroup

	for routeIndex, routeInstance := range r.RouteInstances {
		if routeInstance.Input() == sourceId {
			routeWaitGroup.Go(func() {

				routeFound = true
				routeContext := context.WithValue(ctx, route.SourceContextKey, sourceId)

				payload, err := routeInstance.ProcessPayload(routeContext, payload)
				if err != nil {
					if routeIOErrors == nil {
						routeIOErrors = []route.RouteIOError{}
					}
					r.logger.Error("unable to process input", "route", routeIndex, "source", sourceId, "error", err)
					routeIOErrors = append(routeIOErrors, route.RouteIOError{
						Index:        routeIndex,
						ProcessError: err,
					})
					return
				}

				if payload == nil {
					r.logger.Error("no input after processing", "route", routeIndex, "source", sourceId)
					return
				}

				outputErrors := r.HandleOutput(routeContext, routeInstance.Output(), payload)
				if outputErrors != nil {
					if routeIOErrors == nil {
						routeIOErrors = []route.RouteIOError{}
					}
					routeIOErrors = append(routeIOErrors, route.RouteIOError{
						Index:        routeIndex,
						OutputErrors: outputErrors,
					})
				}
			})

		}
	}
	routeWaitGroup.Wait()
	return routeFound, routeIOErrors
}

func (r *Router) HandleOutput(ctx context.Context, destinationId string, payload any) []error {

	var outputErrors []error
	for _, moduleInstance := range r.ModuleInstances {
		if moduleInstance.Id() == destinationId {
			err := moduleInstance.Output(ctx, payload)
			if err != nil {
				if outputErrors == nil {
					outputErrors = []error{}
				}
				outputErrors = append(outputErrors, err)
				r.logger.Error("unable to route output", "module", moduleInstance.Id(), "error", err)
			}
		}
	}
	return outputErrors
}
