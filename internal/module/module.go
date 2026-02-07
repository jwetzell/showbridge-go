package module

import (
	"context"
	"fmt"
	"log/slog"
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
	Start(context.Context) error
	Stop()
	Output(context.Context, any) error
}

type ModuleRegistration struct {
	Type string `json:"type"`
	New  func(config.ModuleConfig) (Module, error)
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

	if _, ok := ModuleRegistry[string(mod.Type)]; ok {
		panic(fmt.Sprintf("module already registered: %s", mod.Type))
	}
	ModuleRegistry[string(mod.Type)] = mod
}

var (
	moduleRegistryMu sync.RWMutex
	ModuleRegistry   = make(map[string]ModuleRegistration)
)

func CreateLogger(config config.ModuleConfig) *slog.Logger {
	return slog.Default().With("component", "module", "id", config.Id, "type", config.Type)
}
