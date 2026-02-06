package showbridge

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
	"github.com/jwetzell/showbridge-go/internal/route"

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
	tracer         trace.Tracer
}

// TODO(jwetzell): support removing module
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

func NewRouter(config config.Config, tracer trace.Tracer) (*Router, []module.ModuleError, []route.RouteError) {

	router := Router{
		ModuleInstances: make(map[string]module.Module),
		RouteInstances:  []route.Route{},
		logger:          slog.Default().With("component", "router"),
		tracer:          tracer,
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

func (r *Router) Run(ctx context.Context) {
	r.logger.Info("running")
	routerContext, cancel := context.WithCancel(ctx)
	r.Context = routerContext
	r.contextCancel = cancel
	contextWithRouter := context.WithValue(routerContext, route.RouterContextKey, r)

	for _, moduleInstance := range r.ModuleInstances {
		r.moduleWait.Go(func() {
			err := moduleInstance.Run(contextWithRouter)
			if err != nil {
				// TODO(jwetzell): handle module run errors better
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
	r.logger.Info("stopping")
	r.contextCancel()
}

func (r *Router) HandleInput(ctx context.Context, sourceId string, payload any) (bool, []route.RouteIOError) {
	spanCtx, span := r.tracer.Start(ctx, "router.input", trace.WithAttributes(attribute.String("source.id", sourceId)), trace.WithNewRoot())
	defer span.End()
	var routeIOErrors []route.RouteIOError
	routeFound := false

	var routeWaitGroup sync.WaitGroup

	for routeIndex, routeInstance := range r.RouteInstances {
		if routeInstance.Input() == sourceId {
			routeWaitGroup.Go(func() {

				routeFound = true
				routeContext := context.WithValue(spanCtx, route.SourceContextKey, sourceId)

				routeCtx, routeSpan := r.tracer.Start(routeContext, "route", trace.WithAttributes(attribute.Int("route.index", routeIndex), attribute.String("route.input", routeInstance.Input()), attribute.String("route.output", routeInstance.Output())))
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
					r.logger.Error("no input after processing", "route", routeIndex, "source", sourceId)
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
	spanCtx, span := r.tracer.Start(ctx, "router.output", trace.WithAttributes(attribute.String("destination.id", destinationId)))
	defer span.End()

	destinationModule := r.getModule(destinationId)

	if destinationModule == nil {
		err := errors.New("no module found for destination id")
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		r.logger.Error("no module found for destination id", "destinationId", destinationId)
		return err
	}

	moduleOutputCtx, moduleOutputSpan := r.tracer.Start(spanCtx, "module.output", trace.WithAttributes(attribute.String("module.id", destinationModule.Id()), attribute.String("module.type", destinationModule.Type())))
	defer moduleOutputSpan.End()
	err := destinationModule.Output(moduleOutputCtx, payload)
	if err != nil {
		moduleOutputSpan.SetStatus(codes.Error, err.Error())
		moduleOutputSpan.RecordError(err)
		r.logger.Error("module output encountered error", "module", destinationModule.Id(), "error", err)
		return err
	} else {
		moduleOutputSpan.SetStatus(codes.Ok, "module output successful")
	}

	return nil
}
