package showbridge_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/stretchr/testify/require"
)

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
