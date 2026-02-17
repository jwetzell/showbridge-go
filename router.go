package showbridge

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
	"github.com/jwetzell/showbridge-go/internal/route"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Router struct {
	contextCancel   context.CancelFunc
	Context         context.Context
	ModuleInstances map[string]module.Module
	// TODO(jwetzell): change to something easier to lookup
	RouteInstances []route.Route
	moduleWait     sync.WaitGroup
	logger         *slog.Logger
	runningConfig  config.Config
}

func (r *Router) addModule(moduleDecl config.ModuleConfig) error {
	if moduleDecl.Id == "" {
		return errors.New("module id cannot be empty")
	}
	moduleInfo, ok := module.ModuleRegistry[moduleDecl.Type]
	if !ok {
		return errors.New("module type not defined")
	}

	_, ok = r.ModuleInstances[moduleDecl.Id]
	if ok {
		return errors.New("module id already exists")
	}

	moduleInstance, err := moduleInfo.New(moduleDecl)
	if err != nil {
		return err
	}

	r.ModuleInstances[moduleDecl.Id] = moduleInstance
	return nil
}

func (r *Router) removeModule(moduleId string) error {
	err := r.stopModule(moduleId)
	if err != nil {
		return err
	}
	delete(r.ModuleInstances, moduleId)
	return nil
}

func (r *Router) startModule(ctx context.Context, moduleId string) error {
	moduleInstance := r.getModule(moduleId)
	if moduleInstance == nil {
		return errors.New("module id not found")
	}
	r.moduleWait.Go(func() {
		err := moduleInstance.Start(ctx)
		if err != nil {
			// TODO(jwetzell): propagate module run errors better
			r.logger.Error("error encountered running module", "moduleId", moduleId, "error", err)
		}
	})
	return nil
}

func (r *Router) stopModule(moduleId string) error {
	moduleInstance := r.getModule(moduleId)
	if moduleInstance == nil {
		return errors.New("module id not found")
	}
	moduleInstance.Stop()
	return nil
}

// TODO(jwetzell): support removing route
func (r *Router) addRoute(routeDecl config.RouteConfig) error {
	routeInstance, err := route.NewRoute(routeDecl)
	if err != nil {
		return err
	}
	r.RouteInstances = append(r.RouteInstances, routeInstance)
	return nil
}

func (r *Router) getModule(moduleId string) module.Module {
	moduleInstance, ok := r.ModuleInstances[moduleId]
	if !ok {
		return nil
	}
	return moduleInstance
}

func NewRouter(config config.Config) (*Router, []module.ModuleError, []route.RouteError) {

	router := Router{
		ModuleInstances: make(map[string]module.Module),
		RouteInstances:  []route.Route{},
		logger:          slog.Default().With("component", "router"),
		runningConfig:   config,
	}
	router.logger.Debug("creating")

	var moduleErrors []module.ModuleError

	for moduleIndex, moduleDecl := range config.Modules {

		err := router.addModule(moduleDecl)
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

	}

	var routeErrors []route.RouteError
	for routeIndex, routeDecl := range config.Routes {
		err := router.addRoute(routeDecl)
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
	}

	return &router, moduleErrors, routeErrors
}

func (r *Router) Start(ctx context.Context) {
	r.logger.Info("running")
	routerContext, cancel := context.WithCancel(ctx)
	r.Context = routerContext
	r.contextCancel = cancel
	contextWithRouter := context.WithValue(routerContext, route.RouterContextKey, r)

	for moduleId := range r.ModuleInstances {
		// TODO(jwetzell): handle module run errors
		r.startModule(contextWithRouter, moduleId)
	}
	<-r.Context.Done()
	r.logger.Debug("waiting for modules to exit")
	r.moduleWait.Wait()
	r.logger.Info("done")
}

func (r *Router) Stop() {
	r.logger.Info("stopping")
	r.contextCancel()
}

func (r *Router) HandleInput(ctx context.Context, sourceId string, payload any) (bool, []route.RouteIOError) {
	spanCtx, span := otel.Tracer("router").Start(ctx, "input", trace.WithAttributes(attribute.String("source.id", sourceId)), trace.WithNewRoot())
	defer span.End()
	var routeIOErrors []route.RouteIOError
	routeFound := false

	var routeWaitGroup sync.WaitGroup

	for routeIndex, routeInstance := range r.RouteInstances {
		if routeInstance.Input() == sourceId {
			routeWaitGroup.Go(func() {

				routeFound = true
				routeContext := context.WithValue(spanCtx, route.SourceContextKey, sourceId)

				routeCtx, routeSpan := otel.Tracer("router").Start(routeContext, "route", trace.WithAttributes(attribute.Int("route.index", routeIndex), attribute.String("route.input", routeInstance.Input()), attribute.String("route.output", routeInstance.Output())))
				payload, err := routeInstance.ProcessPayload(routeCtx, payload)
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
					r.logger.Debug("no payload after processing, route terminated", "route", routeIndex, "source", sourceId)
					return
				}

				outputError := r.HandleOutput(routeCtx, routeInstance.Output(), payload)
				if outputError != nil {
					if routeIOErrors == nil {
						routeIOErrors = []route.RouteIOError{}
					}
					routeIOErrors = append(routeIOErrors, route.RouteIOError{
						Index:       routeIndex,
						OutputError: outputError,
					})
				}
				routeSpan.End()
			})
		}
	}
	routeWaitGroup.Wait()
	return routeFound, routeIOErrors
}

func (r *Router) HandleOutput(ctx context.Context, destinationId string, payload any) error {
	spanCtx, span := otel.Tracer("router").Start(ctx, "output", trace.WithAttributes(attribute.String("destination.id", destinationId)))
	defer span.End()

	destinationModule := r.getModule(destinationId)

	if destinationModule == nil {
		err := errors.New("no module found for destination id")
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		r.logger.Error("no module found for destination id", "destinationId", destinationId)
		return err
	}

	moduleOutputCtx, moduleOutputSpan := otel.Tracer("module").Start(spanCtx, "output", trace.WithAttributes(attribute.String("module.id", destinationModule.Id()), attribute.String("module.type", destinationModule.Type())))
	defer moduleOutputSpan.End()
	err := destinationModule.Output(moduleOutputCtx, payload)
	if err != nil {
		moduleOutputSpan.SetStatus(codes.Error, err.Error())
		moduleOutputSpan.RecordError(err)
		r.logger.ErrorContext(moduleOutputCtx, "module output encountered error", "module", destinationModule.Id(), "error", err)
		return err
	} else {
		moduleOutputSpan.SetStatus(codes.Ok, "module output successful")
	}

	return nil
}

func (r *Router) RunningConfig() config.Config {
	return r.runningConfig
}
