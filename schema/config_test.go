package schema_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/schema"
)

func TestGoodApplyDefaults(t *testing.T) {

	testCases := []struct {
		name     string
		input    *map[string]any
		expected any
	}{
		{
			name:  "empty config",
			input: &map[string]any{},
			expected: map[string]any{
				"api": map[string]any{
					"enabled": false,
					"port":    float64(8080),
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			err := schema.ApplyDefaults(testCase.input)

			if err != nil {
				t.Fatalf("ApplyDefaults returned error: %v", err)
			}

			if !reflect.DeepEqual(*testCase.input, testCase.expected) {
				t.Fatalf("got %+v, expected %+v", *testCase.input, testCase.expected)
			}
		})
	}

}

func TestGoodValidateConfig(t *testing.T) {

	testCases := []struct {
		name  string
		input map[string]any
	}{
		{
			name:  "empty config",
			input: map[string]any{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			err := schema.ValidateConfig(testCase.input)
			if err != nil {
				t.Fatalf("ValidateConfig returned error: %v", err)
			}
		})
	}

}
