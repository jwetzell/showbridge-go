package showbridge

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/jwetzell/showbridge-go/internal/api"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
	"github.com/jwetzell/showbridge-go/internal/route"
)

// TODO(jwetzell): can/should this be split into different "components"?
type Router struct {
	contextCancel context.CancelFunc
	Context       context.Context
	// TODO(jwetzell): do these need to be guarded against concurrency?
	ModuleInstances map[string]common.Module
	// TODO(jwetzell): change to something easier to lookup
	RouteInstances      []*route.Route
	ConfigChange        chan config.Config
	moduleWait          sync.WaitGroup
	logger              *slog.Logger
	runningConfig       config.Config
	runningConfigMu     sync.RWMutex
	apiServer           *api.ApiServer
	eventDestinations   []common.EventDestination
	eventDestinationsMu sync.Mutex
}

func (r *Router) addModule(moduleDecl config.ModuleConfig) error {
	if moduleDecl.Id == "" {
		return errors.New("module id cannot be empty")
	}
	moduleRegistration, ok := module.GetModuleRegistration(moduleDecl.Type)
	if !ok {
		return errors.New("module type not defined")
	}

	_, ok = r.ModuleInstances[moduleDecl.Id]
	if ok {
		return errors.New("module id already exists")
	}

	moduleInstance, err := moduleRegistration.New(moduleDecl)
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
		err := moduleInstance.Start(ctx, r.HandleInput)
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

func (r *Router) getModule(moduleId string) common.Module {
	moduleInstance, ok := r.ModuleInstances[moduleId]
	if !ok {
		return nil
	}
	return moduleInstance
}

func NewRouter(routerConfig config.Config) (*Router, []config.ModuleError, []config.RouteError) {

	router := Router{
		ModuleInstances: make(map[string]common.Module),
		RouteInstances:  []*route.Route{},
		ConfigChange:    make(chan config.Config, 1),
		logger:          slog.Default().With("component", "router"),
		runningConfig:   routerConfig,
	}
	router.logger.Debug("creating")

	var moduleErrors []config.ModuleError

	for moduleIndex, moduleDecl := range routerConfig.Modules {

		err := router.addModule(moduleDecl)
		if err != nil {
			if moduleErrors == nil {
				moduleErrors = []config.ModuleError{}
			}
			moduleErrors = append(moduleErrors, config.ModuleError{
				Index:  moduleIndex,
				Config: moduleDecl,
				Error:  err.Error(),
			})
			continue
		}

	}

	var routeErrors []config.RouteError
	for routeIndex, routeDecl := range routerConfig.Routes {
		err := router.addRoute(routeDecl)
		if err != nil {
			if routeErrors == nil {
				routeErrors = []config.RouteError{}
			}
			routeErrors = append(routeErrors, config.RouteError{
				Index:  routeIndex,
				Config: routeDecl,
				Error:  err.Error(),
			})
			continue
		}
	}

	apiServer := api.NewApiServer(&router, &router)

	router.apiServer = apiServer

	return &router, moduleErrors, routeErrors
}

func (r *Router) Start(ctx context.Context) {
	r.logger.Info("running")
	routerContext, cancel := context.WithCancel(ctx)
	r.Context = routerContext
	r.contextCancel = cancel
	r.startModules()
	r.apiServer.Start(r.GetRunningConfig().Api)
}

func (r *Router) Stop() {
	r.logger.Info("stopping")
	r.logger.Debug("shutting down api server")
	r.apiServer.Stop()
	r.logger.Debug("stopping modules")
	r.stopModules()
	r.logger.Debug("waiting for modules to exit")
	r.moduleWait.Wait()
	r.logger.Debug("canceling router context")
	r.contextCancel()
	r.logger.Info("done")
}

func (r *Router) HandleInput(ctx context.Context, sourceId string, payload any) (bool, []common.RouteIOError) {
	r.runningConfigMu.RLock()
	defer r.runningConfigMu.RUnlock()

	var routeIOErrors []common.RouteIOError
	var routeFound atomic.Bool

	r.broadcastEvent(common.Event{
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

				routeFound.Store(true)

				_, err := routeInstance.ProcessPayload(ctx, common.WrappedPayload{
					Payload:      payload,
					Source:       sourceId,
					Modules:      r.ModuleInstances,
					InputHandler: r.HandleInput,
					End:          false,
				})
				if err != nil {
					if routeIOErrors == nil {
						routeIOErrors = []common.RouteIOError{}
					}
					r.logger.Error("unable to process input", "route", routeIndex, "source", sourceId, "error", err)
					routeIOErrors = append(routeIOErrors, common.RouteIOError{
						Index:        routeIndex,
						ProcessError: err,
					})
					r.broadcastEvent(common.Event{
						Type: "route",
						Data: map[string]any{
							"index": routeIndex,
						},
						Error: err.Error(),
					})
					return
				}
				r.broadcastEvent(common.Event{
					Type: "route",
					Data: map[string]any{
						"index": routeIndex,
					},
				})
			})
		}
	}
	routeWaitGroup.Wait()
	return routeFound.Load(), routeIOErrors
}

func (r *Router) startModules() {
	for moduleId := range r.ModuleInstances {
		// TODO(jwetzell): handle module run errors
		err := r.startModule(r.Context, moduleId)
		if err != nil {
			r.logger.Error("error starting module", "moduleId", moduleId, "error", err)
		}
	}
}

func (r *Router) stopModules() {
	for moduleId := range r.ModuleInstances {
		// TODO(jwetzell): handle module stop errors?
		err := r.stopModule(moduleId)
		if err != nil {
			r.logger.Error("error stopping module", "moduleId", moduleId, "error", err)
		}
	}
}
