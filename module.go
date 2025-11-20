package showbridge

import (
	"context"
	"fmt"
	"sync"
)

type Module interface {
	Id() string
	Type() string
	RegisterRouter(*Router)
	Run(context.Context) error
	Output(any) error
}

type ModuleConfig struct {
	Id     string         `json:"id"`
	Type   string         `json:"type"`
	Params map[string]any `json:"params"`
}

type ModuleRegistration struct {
	Type string `json:"type"`
	New  func(ModuleConfig) (Module, error)
}

func RegisterModule(mod ModuleRegistration) {

	if mod.Type == "" {
		panic("module type is missing")
	}
	if mod.New == nil {
		panic("missing ModuleInfo.New")
	}

	moduleRegistryMu.Lock()
	defer moduleRegistryMu.Unlock()

	if _, ok := moduleRegistry[string(mod.Type)]; ok {
		panic(fmt.Sprintf("module already registered: %s", mod.Type))
	}
	moduleRegistry[string(mod.Type)] = mod
}

var (
	moduleRegistryMu sync.RWMutex
	moduleRegistry   = make(map[string]ModuleRegistration)
)
