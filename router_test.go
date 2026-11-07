package showbridge_test

import (
	"context"
	"log/slog"
	"reflect"
	"sync"
	"testing"

	"github.com/jwetzell/showbridge-go"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/module"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	module.RegisterModule(module.ModuleRegistration{
		Type: "mock.counter",
		New: func(config config.ModuleConfig) (routing.Module, error) {
			return &MockCounterModule{config: config, logger: slog.Default()}, nil
		},
	})
}

type MockCounterModule struct {
	config      config.ModuleConfig
	ctx         context.Context
	outputCount int
	mu          sync.Mutex
	router      routing.Router
	logger      *slog.Logger
	cancel      context.CancelFunc
}

func (mcm *MockCounterModule) Output(context.Context, any) error {
	mcm.mu.Lock()
	defer mcm.mu.Unlock()
	mcm.outputCount += 1
	return nil
}

func (mcm *MockCounterModule) Start(ctx context.Context, router routing.Router) error {
	mcm.mu.Lock()
	mcm.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	mcm.ctx = moduleContext
	mcm.cancel = cancel
	mcm.mu.Unlock()
	<-mcm.ctx.Done()
	return nil
}

func (mcm *MockCounterModule) Stop() {
	mcm.mu.Lock()
	if mcm.cancel != nil {
		mcm.cancel()
	}
	mcm.mu.Unlock()
}

func TestNewRouter(t *testing.T) {
	routerConfig := config.Config{
		Modules: []config.ModuleConfig{
			{
				Id:   "mock",
				Type: "mock.counter",
			},
		},
		Routes: []config.RouteConfig{},
	}

	router, moduleErrors, routeErrors := showbridge.NewRouter(routerConfig)

	if moduleErrors != nil {
		t.Fatalf("router should not have returned any module errors: %v", moduleErrors)
	}

	if routeErrors != nil {
		t.Fatalf("router should not have returned any route errors: %v", routeErrors)
	}

	if !reflect.DeepEqual(routerConfig, router.GetRunningConfig()) {
		t.Fatalf("router running config did not match expected, got: %v, expected: %v", router.GetRunningConfig(), routerConfig)
	}
}

func TestNewRouterNoModuleId(t *testing.T) {
	routerConfig := config.Config{
		Modules: []config.ModuleConfig{
			{
				Id:   "",
				Type: "mock.counter",
			},
		},
	}

	_, moduleErrors, _ := showbridge.NewRouter(routerConfig)

	if moduleErrors == nil {
		t.Fatalf("router should have returned 'unknown module' module errors")
	}
}

func TestNewRouterUnknownModuleType(t *testing.T) {
	routerConfig := config.Config{
		Modules: []config.ModuleConfig{
			{
				Id:   "mock",
				Type: "asd.fjlkj23oiu4ksldj",
			},
		},
	}

	_, moduleErrors, _ := showbridge.NewRouter(routerConfig)

	if moduleErrors == nil {
		t.Fatalf("router should have returned 'unknown module' module errors")
	}
}

func TestNewRouterDuplicateModuleId(t *testing.T) {
	routerConfig := config.Config{
		Modules: []config.ModuleConfig{
			{
				Id:   "mock",
				Type: "mock.counter",
			},
			{
				Id:   "mock",
				Type: "mock.counter",
			},
		},
	}

	_, moduleErrors, _ := showbridge.NewRouter(routerConfig)

	if moduleErrors == nil {
		t.Fatalf("router should have returned module error")
	}

	if len(moduleErrors) != 1 {
		t.Fatalf("router should have returned exactly 1 module error, got: %d", len(moduleErrors))
	}

	if moduleErrors[0].Error != "module id already exists" {
		t.Fatalf("module error did not match expected, got: %s", moduleErrors[0].Error)
	}
}

func TestNewRouterRouteWithUnknwonProcessor(t *testing.T) {
	routerConfig := config.Config{
		Modules: []config.ModuleConfig{
			{
				Id:   "mock",
				Type: "mock.counter",
			},
		},
		Routes: []config.RouteConfig{
			{
				Input: "mock",
				Processors: []config.ProcessorConfig{
					{
						Type: "asdfasdf",
					},
				},
			},
		},
	}

	_, _, routeErrors := showbridge.NewRouter(routerConfig)

	if routeErrors == nil {
		t.Fatalf("router should have returned a route error")
	}

	if len(routeErrors) != 1 {
		t.Fatalf("router should have returned exactly 1 route error, got: %d", len(routeErrors))
	}

	if routeErrors[0].Error != "problem loading processor registration for processor type: asdfasdf" {
		t.Fatalf("route error did not match expected, got: %s", routeErrors[0].Error)
	}
}
