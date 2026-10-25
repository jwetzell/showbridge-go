package showbridge_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/params"
	"github.com/stretchr/testify/require"
)

func TestGoodNewModule(t *testing.T) {
	testCases := []struct {
		name         string
		moduleConfig config.ModuleConfig
	}{
		{
			name: "good module",
			moduleConfig: config.ModuleConfig{
				Id:   "test-module",
				Type: "time.timer",
				Params: params.Params{
					"duration": 1000,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := showbridge.NewModule(tc.moduleConfig)
			require.NoError(t, err)
		})
	}
}

func TestBadNewModule(t *testing.T) {
	testCases := []struct {
		name         string
		moduleConfig config.ModuleConfig
		errorString  string
	}{
		{
			name: "bad module",
			moduleConfig: config.ModuleConfig{
				Id:   "test-module",
				Type: "bad.type",
			},
			errorString: "problem loading module registration for module type: bad.type",
		},
		{
			name: "empty module id",
			moduleConfig: config.ModuleConfig{
				Id:   "",
				Type: "time.timer",
			},
			errorString: "module id cannot be empty",
		},
		{
			name: "module new error",
			moduleConfig: config.ModuleConfig{
				Id:   "test-module",
				Type: "time.timer",
			},
			errorString: "time.timer duration error: not found",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := showbridge.NewModule(tc.moduleConfig)
			require.EqualError(t, err, tc.errorString)
		})
	}
}
