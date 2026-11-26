package showbridge

import (
	"context"
	"errors"
	"fmt"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/module"
	"github.com/jwetzell/showbridge-go/routing"
)

type moduleManager map[string]routing.Module

func (m moduleManager) GetModule(moduleId string) (routing.Module, bool) {
	moduleInstance, ok := m[moduleId]
	return moduleInstance, ok
}

func NewModule(moduleConfig config.ModuleConfig) (routing.Module, error) {
	if moduleConfig.Id == "" {
		return nil, errors.New("module id cannot be empty")
	}
	moduleRegistration, ok := module.GetModuleRegistration(moduleConfig.Type)
	if !ok {
		return nil, fmt.Errorf("problem loading module registration for module type: %s", moduleConfig.Type)
	}

	moduleInstance, err := moduleRegistration.New(moduleConfig)
	if err != nil {
		return nil, err
	}
	return moduleInstance, nil
}

func (m moduleManager) Get(ctx context.Context, database string, key string) (any, error) {
	moduleInstance, ok := m[database]
	if !ok {
		return nil, fmt.Errorf("module not found for database: %s", database)
	}

	kvModule, ok := moduleInstance.(routing.KeyValueModule)
	if !ok {
		return nil, fmt.Errorf("module for database %s does not implement KeyValueModule", database)
	}
	return kvModule.Get(ctx, key)
}

func (m moduleManager) Set(ctx context.Context, database string, key string, value any) error {
	moduleInstance, ok := m[database]
	if !ok {
		return fmt.Errorf("module not found for database: %s", database)
	}

	kvModule, ok := moduleInstance.(routing.KeyValueModule)
	if !ok {
		return fmt.Errorf("module for database %s does not implement KeyValueModule", database)
	}
	return kvModule.Set(ctx, key, value)
}
