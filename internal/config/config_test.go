package config_test

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
)

func TestGoodStringParamsJSON(t *testing.T) {
	testCases := []struct {
		name       string
		paramsJSON string
		key        string
		expected   string
	}{
		{
			name:       "string param",
			paramsJSON: `{"key": "value"}`,
			key:        "key",
			expected:   "value",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			params := config.Params{}
			err := json.Unmarshal([]byte(testCase.paramsJSON), &params)
			if err != nil {
				t.Fatalf("Failed to unmarshal params JSON: %v", err)
			}
			value, err := params.GetString(testCase.key)
			if err != nil {
				t.Fatalf("GetString returned error: %v", err)
			}
			if value != testCase.expected {
				t.Fatalf("GetString	 got %s, expected %s", value, testCase.expected)
			}
		})
	}
}

func TestGoodIntParamsJSON(t *testing.T) {
	testCases := []struct {
		name       string
		paramsJSON string
		key        string
		expected   int
	}{
		{
			name:       "int param",
			paramsJSON: `{"key": 1}`,
			key:        "key",
			expected:   1,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			params := config.Params{}
			err := json.Unmarshal([]byte(testCase.paramsJSON), &params)
			if err != nil {
				t.Fatalf("Failed to unmarshal params JSON: %v", err)
			}
			value, err := params.GetInt(testCase.key)
			if err != nil {
				t.Fatalf("GetInt returned error: %v", err)
			}
			if value != testCase.expected {
				t.Fatalf("GetInt got %d, expected %d", value, testCase.expected)
			}
		})
	}
}

func TestGoodBoolParamsJSON(t *testing.T) {
	testCases := []struct {
		name       string
		paramsJSON string
		key        string
		expected   bool
	}{
		{
			name:       "bool param",
			paramsJSON: `{"key": true}`,
			key:        "key",
			expected:   true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			params := config.Params{}
			err := json.Unmarshal([]byte(testCase.paramsJSON), &params)
			if err != nil {
				t.Fatalf("Failed to unmarshal params JSON: %v", err)
			}
			value, err := params.GetBool(testCase.key)
			if err != nil {
				t.Fatalf("GetBool returned error: %v", err)
			}
			if value != testCase.expected {
				t.Fatalf("GetBool got %t, expected %t", value, testCase.expected)
			}
		})
	}
}

func TestGoodStringSliceParamsJSON(t *testing.T) {
	testCases := []struct {
		name       string
		paramsJSON string
		key        string
		expected   []string
	}{
		{
			name:       "string array",
			paramsJSON: `{"key": ["value1", "value2"]}`,
			key:        "key",
			expected:   []string{"value1", "value2"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			params := config.Params{}
			err := json.Unmarshal([]byte(testCase.paramsJSON), &params)
			if err != nil {
				t.Fatalf("Failed to unmarshal params JSON: %v", err)
			}
			value, err := params.GetStringSlice(testCase.key)
			if err != nil {
				t.Fatalf("GetStringSlice returned error: %v", err)
			}
			if !slices.Equal(value, testCase.expected) {
				t.Fatalf("GetStringSlice got %v, expected %v", value, testCase.expected)
			}
		})
	}
}

func TestGoodIntSliceParamsJSON(t *testing.T) {
	testCases := []struct {
		name       string
		paramsJSON string
		key        string
		expected   []int
	}{
		{
			name:       "int array",
			paramsJSON: `{"key": [1, 2, 3]}`,
			key:        "key",
			expected:   []int{1, 2, 3},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			params := config.Params{}
			err := json.Unmarshal([]byte(testCase.paramsJSON), &params)
			if err != nil {
				t.Fatalf("Failed to unmarshal params JSON: %v", err)
			}
			value, err := params.GetIntSlice(testCase.key)
			if err != nil {
				t.Fatalf("GetIntSlice returned error: %v", err)
			}
			if !slices.Equal(value, testCase.expected) {
				t.Fatalf("GetIntSlice got %v, expected %v", value, testCase.expected)
			}
		})
	}
}

func TestGoodByteSliceParamsJSON(t *testing.T) {
	testCases := []struct {
		name       string
		paramsJSON string
		key        string
		expected   []byte
	}{
		{
			name:       "byte array",
			paramsJSON: `{"key": [1,2,3,4]}`,
			key:        "key",
			expected:   []byte{1, 2, 3, 4},
		},
		{
			name:       "byte array with floats",
			paramsJSON: `{"key": [1.0,2.0,3.0,4.0]}`,
			key:        "key",
			expected:   []byte{1, 2, 3, 4},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			params := config.Params{}
			err := json.Unmarshal([]byte(testCase.paramsJSON), &params)
			if err != nil {
				t.Fatalf("Failed to unmarshal params JSON: %v", err)
			}
			value, err := params.GetByteSlice(testCase.key)
			if err != nil {
				t.Fatalf("GetByteSlice returned error: %v", err)
			}
			if !slices.Equal(value, testCase.expected) {
				t.Fatalf("GetByteSlice got %v, expected %v", value, testCase.expected)
			}
		})
	}
}
