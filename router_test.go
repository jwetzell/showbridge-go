package showbridge_test

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/jwetzell/showbridge-go"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
	"github.com/jwetzell/showbridge-go/internal/route"
)

type MockCounterModule struct {
	config      config.ModuleConfig
	ctx         context.Context
	outputCount int
	router      route.RouteIO
	logger      *slog.Logger
	cancel      context.CancelFunc
}

func (mcm *MockCounterModule) Id() string {
	return mcm.config.Id
}

func (mcm *MockCounterModule) Output(context.Context, any) error {
	mcm.outputCount += 1
	return nil
}

func (mcm *MockCounterModule) Start(ctx context.Context) error {
	router, ok := ctx.Value(route.RouterContextKey).(route.RouteIO)

	if !ok {
		return fmt.Errorf("mock.counter could not get router from context")
	}
	mcm.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	mcm.ctx = moduleContext
	mcm.cancel = cancel
	<-mcm.ctx.Done()
	return nil
}

func (mcm *MockCounterModule) Type() string {
	return mcm.config.Type
}

func (mcm *MockCounterModule) Stop() {
	mcm.cancel()
}

func init() {
	module.RegisterModule(module.ModuleRegistration{
		Type: "mock.counter",
		New: func(config config.ModuleConfig) (module.Module, error) {
			return &MockCounterModule{config: config, logger: slog.Default()}, nil
		},
	})
}

func TestNewRouter(t *testing.T) {
	routerConfig := config.Config{
		Modules: []config.ModuleConfig{
			{
				Id:   "mock",
				Type: "mock.counter",
			},
		},
	}

	_, moduleErrors, routeErrors := showbridge.NewRouter(routerConfig)

	if moduleErrors != nil {
		t.Fatalf("router should not have returned any module errors: %v", moduleErrors)
	}

	if routeErrors != nil {
		t.Fatalf("router should not have returned any route errors: %v", routeErrors)
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

	if moduleErrors[0].Error.Error() != "module id already exists" {
		t.Fatalf("module error did not match expected, got: %s", moduleErrors[0].Error.Error())
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
				Output: "mock",
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

	if routeErrors[0].Error.Error() != "problem loading processor registration for processor type: asdfasdf" {
		t.Fatalf("route error did not match expected, got: %s", routeErrors[0].Error.Error())
	}
}

func TestRouterInputUnknownDestinationModule(t *testing.T) {
	routerConfig := config.Config{
		Modules: []config.ModuleConfig{
			{
				Id:   "mock",
				Type: "mock.counter",
			},
		},
		Routes: []config.RouteConfig{
			{
				Input:  "mock",
				Output: "test",
			},
		},
	}

	router, moduleErrors, routeErrors := showbridge.NewRouter(routerConfig)

	if moduleErrors != nil {
		t.Fatalf("router should not have returned any module errors: %v", moduleErrors)
	}

	if routeErrors != nil {
		t.Fatalf("router should not have returned any route errors: %v", routeErrors)
	}

	routerRunner := sync.WaitGroup{}

	routerRunner.Go(func() {
		router.Start(t.Context())
		fmt.Println("router stopped")
	})

	time.Sleep(time.Second * 1)

	defer router.Stop()

	_, routingErrors := router.HandleInput(t.Context(), "mock", "test")

	if routingErrors == nil {
		t.Fatalf("router should encounter routing errors when trying to route to an unknown module")
	}

	if len(routingErrors) != 1 {
		t.Fatalf("router should have returned exactly 1 routing error, got: %d", len(routingErrors))
	}

	if routingErrors[0].OutputError.Error() != "no module found for destination id" {
		t.Fatalf("routing output error did not match expected, got: %s", routingErrors[0].OutputError.Error())
	}
}

func TestRouterInputNoMatchingRoute(t *testing.T) {
	routerConfig := config.Config{
		Modules: []config.ModuleConfig{
			{
				Id:   "mock",
				Type: "mock.counter",
			},
		},
		Routes: []config.RouteConfig{
			{
				Input:  "test",
				Output: "mock",
			},
		},
	}

	router, moduleErrors, routeErrors := showbridge.NewRouter(routerConfig)

	if moduleErrors != nil {
		t.Fatalf("router should not have returned any module errors: %v", moduleErrors)
	}

	if routeErrors != nil {
		t.Fatalf("router should not have returned any route errors: %v", routeErrors)
	}

	routerRunner := sync.WaitGroup{}

	routerRunner.Go(func() {
		router.Start(t.Context())
		fmt.Println("router stopped")
	})

	time.Sleep(time.Second * 1)

	defer router.Stop()

	aRouteFound, _ := router.HandleInput(t.Context(), "mock", "test")

	if aRouteFound {
		t.Fatalf("router should not have found a matching route for the input")
	}
}

func TestRouterInputSingleRoute(t *testing.T) {
	routerConfig := config.Config{
		Modules: []config.ModuleConfig{
			{
				Id:   "mock",
				Type: "mock.counter",
			},
		},
		Routes: []config.RouteConfig{
			{
				Input:  "mock",
				Output: "mock",
			},
		},
	}

	router, moduleErrors, routeErrors := showbridge.NewRouter(routerConfig)

	if moduleErrors != nil {
		t.Fatalf("router should not have returned any module errors: %v", moduleErrors)
	}

	if routeErrors != nil {
		t.Fatalf("router should not have returned any route errors: %v", routeErrors)
	}

	routerRunner := sync.WaitGroup{}

	routerRunner.Go(func() {
		router.Start(t.Context())
		fmt.Println("router stopped")
	})

	time.Sleep(time.Second * 1)

	defer router.Stop()

	mockModuleInputCount := 3
	for i := range mockModuleInputCount {
		aRouteFound, routingErrors := router.HandleInput(t.Context(), "mock", fmt.Sprintf("test %d", i))

		if routingErrors != nil {
			t.Fatalf("router should not have encountered routing errors")
		}

		if !aRouteFound {
			t.Fatalf("router should have found a valid route for the input")
		}
	}

	for _, moduleInstance := range router.ModuleInstances {
		if moduleInstance.Id() == "mock" {
			mockModuleInstance, ok := moduleInstance.(*MockCounterModule)
			if !ok {
				t.Fatalf("couldn't get mock module")
			}

			if mockModuleInstance.outputCount != mockModuleInputCount {
				t.Fatalf("mock module output count did not matched expected: %d got: %d", mockModuleInputCount, mockModuleInstance.outputCount)
			}
		}
	}
}

func TestRouterInputMultipleRoutes(t *testing.T) {
	routerConfig := config.Config{
		Modules: []config.ModuleConfig{
			{
				Id:   "mock",
				Type: "mock.counter",
			},
		},
		Routes: []config.RouteConfig{
			{
				Input:  "mock",
				Output: "mock",
			},
			{
				Input:  "mock",
				Output: "mock",
			},
			{
				Input:  "mock",
				Output: "mock",
			},
		},
	}

	router, moduleErrors, routeErrors := showbridge.NewRouter(routerConfig)

	if moduleErrors != nil {
		t.Fatalf("router should not have returned any module errors: %v", moduleErrors)
	}

	if routeErrors != nil {
		t.Fatalf("router should not have returned any route errors: %v", routeErrors)
	}

	routerRunner := sync.WaitGroup{}

	routerRunner.Go(func() {
		router.Start(t.Context())
	})
	time.Sleep(time.Second * 1)

	defer router.Stop()

	mockModuleInputCount := 3
	for i := range mockModuleInputCount {
		aRouteFound, routingErrors := router.HandleInput(t.Context(), "mock", fmt.Sprintf("test %d", i))

		if routingErrors != nil {
			t.Fatalf("router should not have encountered routing errors")
		}

		if !aRouteFound {
			t.Fatalf("router should have found a valid route for the input")
		}
	}

	for _, moduleInstance := range router.ModuleInstances {
		if moduleInstance.Id() == "mock" {
			mockModuleInstance, ok := moduleInstance.(*MockCounterModule)
			if !ok {
				t.Fatalf("couldn't get mock module")
			}

			if mockModuleInstance.outputCount != len(router.RouteInstances)*mockModuleInputCount {
				t.Fatalf("mock module output count did not matched expected: %d got: %d", len(router.RouteInstances)*mockModuleInputCount, mockModuleInstance.outputCount)
			}
			break
		}
	}
}

func TestRouterInputMultipleModules(t *testing.T) {
	routerConfig := config.Config{
		Modules: []config.ModuleConfig{
			{
				Id:   "mock1",
				Type: "mock.counter",
			},
			{
				Id:   "mock2",
				Type: "mock.counter",
			},
		},
		Routes: []config.RouteConfig{
			{
				Input:  "mock1",
				Output: "mock1",
			},
			{
				Input:  "mock2",
				Output: "mock2",
			},
		},
	}

	router, moduleErrors, routeErrors := showbridge.NewRouter(routerConfig)

	if moduleErrors != nil {
		t.Fatalf("router should not have returned any module errors: %v", moduleErrors)
	}

	if routeErrors != nil {
		t.Fatalf("router should not have returned any route errors: %v", routeErrors)
	}

	routerRunner := sync.WaitGroup{}

	routerRunner.Go(func() {
		router.Start(t.Context())
	})

	time.Sleep(time.Second * 1)

	defer router.Stop()

	mock1ModuleInputCount := 3
	for i := range mock1ModuleInputCount {
		aRouteFound, routingErrors := router.HandleInput(t.Context(), "mock1", fmt.Sprintf("test %d", i))

		if routingErrors != nil {
			t.Fatalf("router should not have encountered routing errors")
		}

		if !aRouteFound {
			t.Fatalf("router should have found a valid route for the input")
		}
	}

	mock2ModuleInputCount := 2
	for i := range mock2ModuleInputCount {
		aRouteFound, routingErrors := router.HandleInput(t.Context(), "mock2", fmt.Sprintf("test %d", i))

		if routingErrors != nil {
			t.Fatalf("router should not have encountered routing errors")
		}

		if !aRouteFound {
			t.Fatalf("router should have found a valid route for the input")
		}
	}

	for _, moduleInstance := range router.ModuleInstances {
		if moduleInstance.Id() == "mock1" {
			mockModuleInstance, ok := moduleInstance.(*MockCounterModule)
			if !ok {
				t.Fatalf("couldn't get mock module")
			}

			if mockModuleInstance.outputCount != mock1ModuleInputCount {
				t.Fatalf("mock module output count did not matched expected: %d got: %d", mock1ModuleInputCount, mockModuleInstance.outputCount)
			}
			break
		}
		if moduleInstance.Id() == "mock2" {
			mockModuleInstance, ok := moduleInstance.(*MockCounterModule)
			if !ok {
				t.Fatalf("couldn't get mock module")
			}

			if mockModuleInstance.outputCount != mock2ModuleInputCount {
				t.Fatalf("mock module output count did not matched expected: %d got: %d", mock2ModuleInputCount, mockModuleInstance.outputCount)
			}
			break
		}
	}
}
