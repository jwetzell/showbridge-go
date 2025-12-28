package route_test

import (
	"context"
	"slices"
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
		t.Fatalf("route failed to create: %v", err)
	}

	if testRoute.Input() != routeConfig.Input {
		t.Fatalf("route input does not match expected input")
	}
	if testRoute.Output() != routeConfig.Output {
		t.Fatalf("route output does not match expected output")
	}
}

type MockRouter struct{}

func (mr *MockRouter) HandleInput(sourceId string, payload any) []route.RouteIOError {
	return nil
}

func (mr *MockRouter) HandleOutput(ctx context.Context, destinationId string, payload any) error {
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
		t.Fatalf("route failed to create: %v", err)
	}

	inputData := "test input data"
	payload, err := testRoute.ProcessPayload(context.WithValue(t.Context(), route.RouterContextKey, &MockRouter{}), inputData)
	if err != nil {
		t.Fatalf("route ProcessPayload returned error: %v", err)
	}

	payloadBytes, ok := payload.([]byte)
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
		},
		Output: "output",
	}

	testRoute, err := route.NewRoute(routeConfig)
	if err != nil {
		t.Fatalf("route failed to create: %v", err)
	}

	inputData := "test input data"
	_, err = testRoute.ProcessPayload(context.WithValue(t.Context(), route.RouterContextKey, &MockRouter{}), inputData)
	if err == nil {
		t.Fatalf("route HandleOutput did not return error for bad processor")
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
		t.Fatalf("route failed to create: %v", err)
		return
	}

	payload, err := testRoute.ProcessPayload(context.WithValue(t.Context(), route.RouterContextKey, &MockRouter{}), nil)
	if err != nil {
		t.Fatalf("route ProcessPayload returned error: %v", err)
	}
	if payload != nil {
		t.Fatalf("route returned the wrong payload")
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
		t.Fatalf("route failed to create: %v", err)
	}

	payload, err := testRoute.ProcessPayload(context.WithValue(t.Context(), route.RouterContextKey, &MockRouter{}), "test")
	if err != nil {
		t.Fatalf("route HandleOutput returned error for nil payload: %v", err)
	}

	if payload != nil {
		t.Fatalf("route returned the wrong payload")
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
		t.Fatalf("route error expected when creating route with an unknown processor, got nil")
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
		t.Fatalf("route error expected creating route with bad processor, got nil")
	}
}
