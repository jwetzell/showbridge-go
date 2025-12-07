package showbridge

import (
	"fmt"
	"sync"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type ModuleError struct {
	Index  int
	Config config.ModuleConfig
	Error  error
}

type Module interface {
	Id() string
	Type() string
	Run() error
	Output(any) error
}

type ModuleRegistration struct {
	Type string `json:"type"`
	New  func(config.ModuleConfig, *Router) (Module, error)
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
