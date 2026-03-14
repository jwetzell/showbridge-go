package showbridge

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"reflect"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
	"github.com/jwetzell/showbridge-go/internal/route"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Router struct {
	contextCancel context.CancelFunc
	Context       context.Context
	// TODO(jwetzell): do these need to be guarded against concurrency?
	ModuleInstances map[string]module.Module
	// TODO(jwetzell): change to something easier to lookup
	RouteInstances    []*route.Route
	ConfigChange      chan config.Config
	moduleWait        sync.WaitGroup
	logger            *slog.Logger
	runningConfig     config.Config
	runningConfigMu   sync.RWMutex
	wsConns           []*websocket.Conn
	wsConnsMu         sync.Mutex
	apiServer         *http.Server
	apiServerMu       sync.Mutex
	apiServerShutdown context.CancelFunc
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

func NewRouter(routerConfig config.Config) (*Router, []module.ModuleError, []route.RouteError) {

	router := Router{
		ModuleInstances: make(map[string]module.Module),
		RouteInstances:  []*route.Route{},
		ConfigChange:    make(chan config.Config, 1),
		logger:          slog.Default().With("component", "router"),
		runningConfig:   routerConfig,
	}
	router.logger.Debug("creating")

	var moduleErrors []module.ModuleError

	for moduleIndex, moduleDecl := range routerConfig.Modules {

		err := router.addModule(moduleDecl)
		if err != nil {
			if moduleErrors == nil {
				moduleErrors = []module.ModuleError{}
			}
			moduleErrors = append(moduleErrors, module.ModuleError{
				Index:  moduleIndex,
				Config: moduleDecl,
				Error:  err.Error(),
			})
			continue
		}

	}

	var routeErrors []route.RouteError
	for routeIndex, routeDecl := range routerConfig.Routes {
		err := router.addRoute(routeDecl)
		if err != nil {
			if routeErrors == nil {
				routeErrors = []route.RouteError{}
			}
			routeErrors = append(routeErrors, route.RouteError{
				Index:  routeIndex,
				Config: routeDecl,
				Error:  err.Error(),
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
	r.startModules()
	r.startAPIServer(r.runningConfig.Api)
	<-r.Context.Done()
	r.logger.Debug("shutting down api server")
	r.stopAPIServer()
	r.logger.Debug("waiting for modules to exit")
	r.moduleWait.Wait()
	r.logger.Info("done")
}

func (r *Router) Stop() {
	r.logger.Info("stopping")
	r.contextCancel()
}

func (r *Router) HandleInput(ctx context.Context, sourceId string, payload any) (bool, []common.RouteIOError) {
	r.runningConfigMu.RLock()
	defer r.runningConfigMu.RUnlock()

	spanCtx, span := otel.Tracer("router").Start(ctx, "input", trace.WithAttributes(attribute.String("source.id", sourceId)), trace.WithNewRoot())
	defer span.End()
	var routeIOErrors []common.RouteIOError
	routeFound := false

	r.broadcastEvent(Event{
		Type: "input",
		Data: map[string]any{
			"source": sourceId,
		},
	})

	var routeWaitGroup sync.WaitGroup

	for routeIndex, routeInstance := range r.RouteInstances {
		if routeInstance == nil {
			r.logger.Error("nil route instance found", "routeIndex", routeIndex)
			continue
		}
		if routeInstance.Input() == sourceId {
			routeWaitGroup.Go(func() {

				routeFound = true
				routeContext := context.WithValue(spanCtx, common.SourceContextKey, sourceId)
				routeContext = context.WithValue(routeContext, common.RouterContextKey, r)
				routeContext = context.WithValue(routeContext, common.ModulesContextKey, r.ModuleInstances)

				routeCtx, routeSpan := otel.Tracer("router").Start(routeContext, "route", trace.WithAttributes(attribute.Int("route.index", routeIndex), attribute.String("route.input", routeInstance.Input())))
				_, err := routeInstance.ProcessPayload(routeCtx, payload)
				if err != nil {
					if routeIOErrors == nil {
						routeIOErrors = []common.RouteIOError{}
					}
					r.logger.Error("unable to process input", "route", routeIndex, "source", sourceId, "error", err)
					routeIOErrors = append(routeIOErrors, common.RouteIOError{
						Index:        routeIndex,
						ProcessError: err,
					})
					r.broadcastEvent(Event{
						Type: "route",
						Data: map[string]any{
							"index": routeIndex,
						},
						Error: err.Error(),
					})
					return
				}
				r.broadcastEvent(Event{
					Type: "route",
					Data: map[string]any{
						"index": routeIndex,
					},
				})
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
	outputEvent := Event{
		Type: "output",
		Data: map[string]any{
			"destination": destinationId,
		},
	}
	destinationModule := r.getModule(destinationId)

	if destinationModule == nil {
		err := errors.New("no module found for destination id")
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		r.logger.Error("no module found for destination id", "destinationId", destinationId)
		outputEvent.Error = err.Error()
		r.broadcastEvent(outputEvent)
		return err
	}

	moduleOutputCtx, moduleOutputSpan := otel.Tracer("module").Start(spanCtx, "output", trace.WithAttributes(attribute.String("module.id", destinationModule.Id()), attribute.String("module.type", destinationModule.Type())))
	defer moduleOutputSpan.End()
	err := destinationModule.Output(moduleOutputCtx, payload)
	if err != nil {
		moduleOutputSpan.SetStatus(codes.Error, err.Error())
		moduleOutputSpan.RecordError(err)
		r.logger.ErrorContext(moduleOutputCtx, "module output encountered error", "module", destinationModule.Id(), "error", err)
		outputEvent.Error = err.Error()
		r.broadcastEvent(outputEvent)
		return err
	} else {
		moduleOutputSpan.SetStatus(codes.Ok, "module output successful")
	}
	r.broadcastEvent(outputEvent)
	return nil
}

func (r *Router) startModules() {
	contextWithRouter := context.WithValue(r.Context, common.RouterContextKey, r)

	for moduleId := range r.ModuleInstances {
		// TODO(jwetzell): handle module run errors
		err := r.startModule(contextWithRouter, moduleId)
		if err != nil {
			r.logger.Error("error starting module", "moduleId", moduleId, "error", err)
		}
	}
}

func (r *Router) RunningConfig() config.Config {
	r.runningConfigMu.Lock()
	defer r.runningConfigMu.Unlock()
	return r.runningConfig
}

func (r *Router) UpdateConfig(newConfig config.Config) ([]module.ModuleError, []route.RouteError) {
	r.runningConfigMu.Lock()
	defer r.runningConfigMu.Unlock()
	oldConfig := r.runningConfig
	r.logger.Debug("received config update", "oldConfig", oldConfig, "newConfig", newConfig)

	if !reflect.DeepEqual(oldConfig.Api, newConfig.Api) {
		r.logger.Info("applying new API config")
		r.stopAPIServer()
		r.startAPIServer(newConfig.Api)
		r.runningConfig.Api = newConfig.Api
	}

	// TODO(jwetzell): handle config update errors better
	for _, moduleInstance := range r.ModuleInstances {
		moduleInstance.Stop()
	}
	r.logger.Debug("waiting for modules to exit")
	r.moduleWait.Wait()

	r.ModuleInstances = make(map[string]module.Module)
	r.RouteInstances = []*route.Route{}

	var moduleErrors []module.ModuleError

	for moduleIndex, moduleDecl := range newConfig.Modules {

		err := r.addModule(moduleDecl)
		if err != nil {
			if moduleErrors == nil {
				moduleErrors = []module.ModuleError{}
			}
			moduleErrors = append(moduleErrors, module.ModuleError{
				Index:  moduleIndex,
				Config: moduleDecl,
				Error:  err.Error(),
			})
			continue
		}

	}

	var routeErrors []route.RouteError
	for routeIndex, routeDecl := range newConfig.Routes {
		err := r.addRoute(routeDecl)
		if err != nil {
			if routeErrors == nil {
				routeErrors = []route.RouteError{}
			}
			routeErrors = append(routeErrors, route.RouteError{
				Index:  routeIndex,
				Config: routeDecl,
				Error:  err.Error(),
			})
			continue
		}
	}
	r.runningConfig = newConfig
	r.startModules()

	return moduleErrors, routeErrors
}
