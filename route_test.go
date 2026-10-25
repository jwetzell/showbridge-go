package showbridge_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/params"
	"github.com/jwetzell/showbridge-go/routing"
	"github.com/stretchr/testify/require"
)

func TestGoodNewRoute(t *testing.T) {
	testCases := []struct {
		name  string
		route config.RouteConfig
	}{
		{
			name: "basic route",
			route: config.RouteConfig{
				Input: "source1",
			},
		},
		{
			name: "route with processor",
			route: config.RouteConfig{
				Input: "source1",
				Processors: []config.ProcessorConfig{
					{
						Type: "float.parse",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := showbridge.NewRoute(tc.route)
			require.NoError(t, err)
		})
	}
}

func TestBadNewRoute(t *testing.T) {
	testCases := []struct {
		name          string
		route         config.RouteConfig
		expectedError string
	}{
		{
			name: "bad processor type",
			route: config.RouteConfig{
				Input: "source1",
				Processors: []config.ProcessorConfig{
					{
						Type: "dsafsdf;lkj",
					},
				},
			},
			expectedError: "problem loading processor registration for processor type: dsafsdf;lkj",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := showbridge.NewRoute(tc.route)
			require.EqualError(t, err, tc.expectedError)
		})
	}
}

func TestRouteMatches(t *testing.T) {

	testCases := []struct {
		name     string
		source   string
		route    config.RouteConfig
		expected bool
	}{
		{
			name:   "string match",
			source: "source1",
			route: config.RouteConfig{
				Input: "source1",
			},
			expected: true,
		},
		{
			name:   "string mismatch",
			source: "source1",
			route: config.RouteConfig{
				Input: "source2",
			},
			expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			route, err := showbridge.NewRoute(tc.route)

			require.NoError(t, err)

			require.Equal(t, tc.expected, route.Matches(tc.source))
		})
	}
}

func TestGoodRouteProcessPayload(t *testing.T) {
	testCases := []struct {
		name           string
		route          config.RouteConfig
		wrappedPayload routing.WrappedPayload
		moduleManager  routing.ModuleManager
		inputHandler   routing.InputHandler
	}{
		{
			name: "no processors",
			route: config.RouteConfig{
				Input:      "source1",
				Processors: []config.ProcessorConfig{},
			},
			wrappedPayload: routing.WrappedPayload{
				Payload: "test",
				Source:  "source1",
				End:     false,
			},
			moduleManager: nil,
			inputHandler:  nil,
		},
		{
			name: "single processor run to end",
			route: config.RouteConfig{
				Input: "source1",
				Processors: []config.ProcessorConfig{
					{
						Type:   "float.parse",
						Params: params.Params{},
					},
				},
			},
			wrappedPayload: routing.WrappedPayload{
				Payload: "123.45",
				Source:  "source1",
				End:     false,
			},
			moduleManager: nil,
			inputHandler:  nil,
		},
		{
			name: "single processor stop early",
			route: config.RouteConfig{
				Input: "source1",
				Processors: []config.ProcessorConfig{
					{
						Type: "filter.expr",
						Params: params.Params{
							"expression": "Payload != \"123.45\"",
						},
					},
					{
						Type:   "float.parse",
						Params: params.Params{},
					},
				},
			},
			wrappedPayload: routing.WrappedPayload{
				Payload: "123.45",
				Source:  "source1",
				End:     false,
			},
			moduleManager: nil,
			inputHandler:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			route, err := showbridge.NewRoute(tc.route)
			require.NoError(t, err)

			err = route.ProcessPayload(
				t.Context(),
				tc.wrappedPayload,
				tc.moduleManager,
				tc.inputHandler,
			)
			require.NoError(t, err)
		})
	}
}

func TestBadRouteProcessPayload(t *testing.T) {
	testCases := []struct {
		name           string
		route          config.RouteConfig
		wrappedPayload routing.WrappedPayload
		moduleManager  routing.ModuleManager
		inputHandler   routing.InputHandler
		expectedError  string
	}{
		{
			name: "processor error",
			route: config.RouteConfig{
				Processors: []config.ProcessorConfig{
					{
						Type:   "float.parse",
						Params: params.Params{},
					},
				},
			},
			wrappedPayload: routing.WrappedPayload{
				Payload: "test",
				Source:  "source1",
				End:     false,
			},
			moduleManager: nil,
			inputHandler:  nil,
			expectedError: "processor[0] error: strconv.ParseFloat: parsing \"test\": invalid syntax",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			route, err := showbridge.NewRoute(tc.route)
			require.NoError(t, err)

			err = route.ProcessPayload(
				t.Context(),
				tc.wrappedPayload,
				tc.moduleManager,
				tc.inputHandler,
			)
			require.EqualError(t, err, tc.expectedError)
		})
	}
}
