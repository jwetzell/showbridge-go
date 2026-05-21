package route_test

import (
	"context"
	"slices"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/route"
	"github.com/jwetzell/showbridge-go/internal/test"
)

func TestRouteCreate(t *testing.T) {
	routeConfig := config.RouteConfig{
		Input: "input",
	}

	testRoute, err := route.NewRoute(routeConfig)
	if err != nil {
		t.Fatalf("route failed to create: %v", err)
	}

	if testRoute.Input() != routeConfig.Input {
		t.Fatalf("route input does not match expected input")
	}
}

type MockRouter struct{}

func (mr *MockRouter) HandleInput(ctx context.Context, sourceId string, payload any) (bool, []common.RouteIOError) {
	return false, []common.RouteIOError{}
}

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
	testRouter := test.GetNewTestRouter()
	payload, err := testRoute.ProcessPayload(t.Context(), common.WrappedPayload{
		InputHandler:  testRouter.HandleInput,
		Modules: map[string]common.Module{"output": &test.TestOutputModule{}},
		Payload: inputData,
	})
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
	testRouter := test.GetNewTestRouter()
	_, err = testRoute.ProcessPayload(t.Context(), common.WrappedPayload{
		InputHandler:  testRouter.HandleInput,
		Modules: map[string]common.Module{"output": &test.TestOutputModule{}},
		Payload: inputData,
	})
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

	testRouter := test.GetNewTestRouter()
	payload, err := testRoute.ProcessPayload(t.Context(), common.WrappedPayload{
		InputHandler:  testRouter.HandleInput,
		Modules: map[string]common.Module{"output": &test.TestOutputModule{}},
		Payload: nil,
	})
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

	testRouter := test.GetNewTestRouter()
	_, err = testRoute.ProcessPayload(t.Context(), common.WrappedPayload{
		InputHandler:  testRouter.HandleInput,
		Modules: map[string]common.Module{"output": &test.TestOutputModule{}},
		Payload: "test",
	})
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
