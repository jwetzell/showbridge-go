package showbridge

import (
	"errors"
	"reflect"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/route"
)

func (r *Router) GetRunningConfig() config.Config {
	r.runningConfigMu.RLock()
	defer r.runningConfigMu.RUnlock()
	return r.runningConfig
}

func (r *Router) UpdateConfig(newConfig config.Config, triggerChangeChan bool) (error, []config.ModuleError, []config.RouteError) {
	if !r.runningConfigMu.TryLock() {
		return errors.New("config update in progress"), nil, nil
	}
	defer r.runningConfigMu.Unlock()
	oldConfig := r.runningConfig
	r.logger.Debug("received config update", "oldConfig", oldConfig, "newConfig", newConfig)

	if !reflect.DeepEqual(oldConfig.Api, newConfig.Api) {
		r.logger.Info("applying new API config")
		r.apiServer.Stop()
		r.apiServer.Start(newConfig.Api)
		r.runningConfig.Api = newConfig.Api
	}

	// TODO(jwetzell): handle config update errors better
	for _, moduleInstance := range r.ModuleInstances {
		moduleInstance.Stop()
	}
	r.logger.Debug("waiting for modules to exit")
	r.moduleWait.Wait()

	r.ModuleInstances = make(map[string]common.Module)
	r.RouteInstances = []*route.Route{}

	var moduleErrors []config.ModuleError

	for moduleIndex, moduleDecl := range newConfig.Modules {

		err := r.addModule(moduleDecl)
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
	for routeIndex, routeDecl := range newConfig.Routes {
		err := r.addRoute(routeDecl)
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
	r.runningConfig = newConfig
	r.startModules()

	if triggerChangeChan {
		r.ConfigChange <- newConfig
	}

	return nil, moduleErrors, routeErrors
}
