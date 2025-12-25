package route_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/route"
)

func TestRouteCreate(t *testing.T) {
	routeConfig := config.RouteConfig{
		Input:  "input",
		Output: "output",
	}

	testRoute, err := route.NewRoute(routeConfig)
	if err != nil {
		t.Errorf("Failed to create route: %v", err)
		return
	}

	if testRoute.Input() != routeConfig.Input {
		t.Errorf("route input does not match expected input")
	}
	if testRoute.Output() != routeConfig.Output {
		t.Errorf("route output does not match expected output")
	}
}

type MockRouter struct{}

func (mr *MockRouter) HandleInput(sourceId string, payload any) []route.RouteIOError {
	return nil
}

func (mr *MockRouter) HandleOutput(sourceId string, destinationId string, payload any) error {
	return nil
}

func TestGoodRouteHandleInput(t *testing.T) {
	routeConfig := config.RouteConfig{
		Input: "input",
		Processors: []config.ProcessorConfig{
			{Type: "string.encode"},
		},
		Output: "output",
	}

	testRoute, err := route.NewRoute(routeConfig)
	if err != nil {
		t.Errorf("Failed to create route: %v", err)
		return
	}

	inputData := "test input data"
	err = testRoute.HandleInput(t.Context(), "input", inputData, &MockRouter{})
	if err != nil {
		t.Errorf("Route HandleOutput returned error: %v", err)
	}
}

func TestRouteHandleInputWithProcessorError(t *testing.T) {
	routeConfig := config.RouteConfig{
		Input: "input",
		Processors: []config.ProcessorConfig{
			{Type: "string.create", Params: map[string]any{"template": "{{.invalid}}}"}},
		},
		Output: "output",
	}

	testRoute, err := route.NewRoute(routeConfig)
	if err != nil {
		t.Errorf("Failed to create route: %v", err)
		return
	}

	inputData := "test input data"
	err = testRoute.HandleInput(t.Context(), "input", inputData, &MockRouter{})
	if err == nil {
		t.Errorf("Route HandleOutput did not return error for bad processor")
	}
}

func TestRouteHandleNilPayload(t *testing.T) {
	routeConfig := config.RouteConfig{
		Input:      "input",
		Processors: []config.ProcessorConfig{},
		Output:     "output",
	}

	testRoute, err := route.NewRoute(routeConfig)
	if err != nil {
		t.Errorf("Failed to create route: %v", err)
		return
	}

	err = testRoute.HandleInput(t.Context(), "input", nil, &MockRouter{})
	if err != nil {
		t.Errorf("Route HandleOutput returned error for nil payload: %v", err)
	}
}

func TestRouteHandleNilPayloadFromProcessor(t *testing.T) {
	routeConfig := config.RouteConfig{
		Input: "input",
		Processors: []config.ProcessorConfig{
			{Type: "script.js", Params: map[string]any{"program": "payload = undefined"}},
		},
		Output: "output",
	}

	testRoute, err := route.NewRoute(routeConfig)
	if err != nil {
		t.Errorf("Failed to create route: %v", err)
		return
	}

	err = testRoute.HandleInput(t.Context(), "input", nil, &MockRouter{})
	if err != nil {
		t.Errorf("Route HandleOutput returned error for nil payload: %v", err)
	}
}

func TestRouteUnknownProcessor(t *testing.T) {
	routeConfig := config.RouteConfig{
		Input: "input",
		Processors: []config.ProcessorConfig{
			{Type: "asdfasdflkjalkj"},
		},
		Output: "output",
	}

	_, err := route.NewRoute(routeConfig)
	if err == nil {
		t.Errorf("Expected error when creating route with unknown processor, got nil")
		return
	}
}

func TestRouteBadProcessorConfig(t *testing.T) {
	routeConfig := config.RouteConfig{
		Input: "input",
		Processors: []config.ProcessorConfig{
			{Type: "string.create", Params: map[string]any{}},
		},
		Output: "output",
	}

	_, err := route.NewRoute(routeConfig)
	if err == nil {
		t.Errorf("Expected error when creating route with bad processor, got nil")
		return
	}
}
