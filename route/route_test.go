package route_test

import (
	"context"
	"slices"
	"testing"

	"github.com/jwetzell/showbridge-go/common"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/payload"
	"github.com/jwetzell/showbridge-go/route"
)

type MockRouter struct{}

func (mr *MockRouter) HandleInput(ctx context.Context, sourceId string, payload any) (bool, []common.RouteIOError) {
	return false, []common.RouteIOError{}
}

type mockModuleManager map[string]common.Module

func (m mockModuleManager) GetModule(moduleId string) (common.Module, bool) {
	moduleInstance, ok := m[moduleId]
	return moduleInstance, ok
}

func NewTestOutputModule(id string) *TestOutputModule {
	return &TestOutputModule{
		id: id,
	}
}

type TestOutputModule struct {
	id string
}

func (m *TestOutputModule) Start(ctx context.Context, inputHandler common.InputHandler) error {
	<-ctx.Done()
	return nil
}

func (m *TestOutputModule) Output(ctx context.Context, payload any) error {
	return nil
}

func (m *TestOutputModule) Stop() {}

func TestGoodRouteHandleInput(t *testing.T) {
	routeConfig := config.RouteConfig{
		Input: "input",
		Processors: []config.ProcessorConfig{
			{Type: "string.encode"},
			{
				Type: "module.output",
				Params: config.Params{
					"module": "output",
				},
			},
		},
	}

	testRoute, err := route.NewRoute(routeConfig)
	if err != nil {
		t.Fatalf("route failed to create: %v", err)
	}

	inputData := "test input data"
	testRouter := &MockRouter{}
	payload, err := testRoute.ProcessPayload(t.Context(), payload.WrappedPayload{
		Payload: inputData,
	}, mockModuleManager{"output": &TestOutputModule{}}, testRouter.HandleInput)
	if err != nil {
		t.Fatalf("route ProcessPayload returned error: %v", err)
	}

	payloadBytes, ok := common.GetAnyAsByteSlice(payload)
	if !ok {
		t.Fatalf("payload should be []byte got %T", payload)
	}

	if !slices.Equal([]byte(inputData), payloadBytes) {
		t.Fatalf("route returned the wrong payload. expected: %+v got %+v", inputData, payloadBytes)
	}
}

func TestRouteHandleInputWithProcessorError(t *testing.T) {
	routeConfig := config.RouteConfig{
		Input: "input",
		Processors: []config.ProcessorConfig{
			{Type: "string.create", Params: map[string]any{"template": "{{.invalid}}}"}},
			{
				Type: "module.output",
				Params: config.Params{
					"module": "output",
				},
			},
		},
	}

	testRoute, err := route.NewRoute(routeConfig)
	if err != nil {
		t.Fatalf("route failed to create: %v", err)
	}

	inputData := "test input data"
	testRouter := &MockRouter{}
	_, err = testRoute.ProcessPayload(t.Context(), payload.WrappedPayload{
		Payload: inputData,
	}, mockModuleManager{"output": &TestOutputModule{}}, testRouter.HandleInput)
	if err == nil {
		t.Fatalf("route did not return error for bad processor")
	}
}

func TestRouteHandleNilPayload(t *testing.T) {
	routeConfig := config.RouteConfig{
		Input: "input",
		Processors: []config.ProcessorConfig{
			{
				Type: "module.output",
				Params: config.Params{
					"module": "output",
				},
			},
		},
	}

	testRoute, err := route.NewRoute(routeConfig)
	if err != nil {
		t.Fatalf("route failed to create: %v", err)
		return
	}

	testRouter := &MockRouter{}
	payload, err := testRoute.ProcessPayload(t.Context(), payload.WrappedPayload{
		Payload: nil,
	}, mockModuleManager{"output": &TestOutputModule{}}, testRouter.HandleInput)
	if err != nil {
		t.Fatalf("route ProcessPayload returned error: %v", err)
	}
	if payload != nil {
		t.Fatalf("route returned the wrong payload: expected nil got %+v (%T)", payload, payload)
	}
}

func TestRouteHandleNilPayloadFromProcessor(t *testing.T) {
	routeConfig := config.RouteConfig{
		Input: "input",
		Processors: []config.ProcessorConfig{
			{Type: "script.js", Params: map[string]any{"program": "payload = undefined"}},
			{
				Type: "module.output",
				Params: config.Params{
					"module": "output",
				},
			},
		},
	}

	testRoute, err := route.NewRoute(routeConfig)
	if err != nil {
		t.Fatalf("route failed to create: %v", err)
	}

	testRouter := &MockRouter{}
	_, err = testRoute.ProcessPayload(t.Context(), payload.WrappedPayload{
		Payload: "test",
	}, mockModuleManager{"output": &TestOutputModule{}}, testRouter.HandleInput)
	if err != nil {
		t.Fatalf("route returned error for nil payload: %v", err)
	}
}

func TestRouteUnknownProcessor(t *testing.T) {
	routeConfig := config.RouteConfig{
		Input: "input",
		Processors: []config.ProcessorConfig{
			{Type: "asdfasdflkjalkj"},
		},
	}

	_, err := route.NewRoute(routeConfig)
	if err == nil {
		t.Fatalf("route error expected when creating route with an unknown processor, got nil")
	}
}

func TestRouteBadProcessorConfig(t *testing.T) {
	routeConfig := config.RouteConfig{
		Input: "input",
		Processors: []config.ProcessorConfig{
			{Type: "string.create", Params: map[string]any{}},
		},
	}

	_, err := route.NewRoute(routeConfig)
	if err == nil {
		t.Fatalf("route error expected creating route with bad processor, got nil")
	}
}
