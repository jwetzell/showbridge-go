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
	ModuleInstances []module.Module
	RouteInstances  []route.Route
	moduleWait      sync.WaitGroup
	logger          *slog.Logger
	tracer          trace.Tracer
}

func NewRouter(config config.Config, tracer trace.Tracer) (*Router, []module.ModuleError, []route.RouteError) {

	router := Router{
		ModuleInstances: []module.Module{},
		RouteInstances:  []route.Route{},
		logger:          slog.Default().With("component", "router"),
		tracer:          tracer,
	}
	router.logger.Debug("creating")

	var moduleErrors []module.ModuleError

	for moduleIndex, moduleDecl := range config.Modules {

		if moduleDecl.Id == "" {
			if moduleErrors == nil {
				moduleErrors = []module.ModuleError{}
			}
			moduleErrors = append(moduleErrors, module.ModuleError{
				Index:  moduleIndex,
				Config: moduleDecl,
				Error:  errors.New("module id cannot be empty"),
			})
			continue
		}

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
			moduleInstance, err := moduleInfo.New(moduleDecl)
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

				routeSpanCtx, routeSpan := r.tracer.Start(routeContext, "route.input", trace.WithAttributes(attribute.Int("route.index", routeIndex)))
				defer routeSpan.End()
				routeProcessCtx, routeSpan := r.tracer.Start(routeSpanCtx, "route.process")
				payload, err := routeInstance.ProcessPayload(routeProcessCtx, payload)
				if err != nil {
					if routeIOErrors == nil {
						routeIOErrors = []route.RouteIOError{}
					}
					r.logger.Error("unable to process input", "route", routeIndex, "source", sourceId, "error", err)
					routeIOErrors = append(routeIOErrors, route.RouteIOError{
						Index:        routeIndex,
						ProcessError: err,
					})
					routeSpan.SetStatus(codes.Error, err.Error())
					routeSpan.RecordError(err)
					routeSpan.End()
					return
				} else {
					routeSpan.SetStatus(codes.Ok, "route processing successful")
					routeSpan.End()
				}

				if payload == nil {
					r.logger.Error("no input after processing", "route", routeIndex, "source", sourceId)
					return
				}

				routeOutputCtx, routeOutputSpan := r.tracer.Start(routeSpanCtx, "route.output", trace.WithAttributes(attribute.String("destination.id", routeInstance.Output())))
				outputErrors := r.HandleOutput(routeOutputCtx, routeInstance.Output(), payload)
				if outputErrors != nil {
					if routeIOErrors == nil {
						routeIOErrors = []route.RouteIOError{}
					}
					routeIOErrors = append(routeIOErrors, route.RouteIOError{
						Index:        routeIndex,
						OutputErrors: outputErrors,
					})
					routeOutputSpan.SetStatus(codes.Error, "route output error")
					for _, outputError := range outputErrors {
						routeOutputSpan.RecordError(outputError)
					}
				} else {
					routeOutputSpan.SetStatus(codes.Ok, "route output successful")
				}
				routeOutputSpan.End()
			})
		}
	}
	routeWaitGroup.Wait()
	return routeFound, routeIOErrors
}

func (r *Router) HandleOutput(ctx context.Context, destinationId string, payload any) []error {
	spanCtx, span := r.tracer.Start(ctx, "router.output", trace.WithAttributes(attribute.String("destination.id", destinationId)))
	defer span.End()
	var outputErrors []error
	for _, moduleInstance := range r.ModuleInstances {
		if moduleInstance.Id() == destinationId {
			moduleSpanCtx, moduleSpan := r.tracer.Start(spanCtx, "module.output", trace.WithAttributes(attribute.String("module.id", moduleInstance.Id()), attribute.String("module.type", moduleInstance.Type())))
			err := moduleInstance.Output(moduleSpanCtx, payload)
			if err != nil {
				if outputErrors == nil {
					outputErrors = []error{}
				}
				outputErrors = append(outputErrors, err)
				moduleSpan.SetStatus(codes.Error, err.Error())
				moduleSpan.RecordError(err)
				r.logger.Error("module output encountered error", "module", moduleInstance.Id(), "error", err)
			} else {
				moduleSpan.SetStatus(codes.Ok, "module output successful")
			}
			moduleSpan.End()
		}
	}

	if outputErrors != nil {
		span.SetStatus(codes.Error, "router output error")
		for _, outputError := range outputErrors {
			span.RecordError(outputError)
		}
	} else {
		span.SetStatus(codes.Ok, "router output successful")
	}
	return outputErrors
}
