package config_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/stretchr/testify/require"
)

func TestCleanConfig(t *testing.T) {
	testCases := []struct {
		name     string
		input    config.Config
		expected config.Config
	}{
		{
			name:  "empty config",
			input: config.Config{},
			expected: config.Config{
				Modules: []config.ModuleConfig{},
				Routes:  []config.RouteConfig{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := config.CleanConfig(tc.input)
			require.Equal(t, tc.expected, actual)
		})
	}
}
