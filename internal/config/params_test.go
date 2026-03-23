package config_test

import (
	"encoding/json"
	"errors"
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

func TestBadStringParamsJSON(t *testing.T) {
	testCases := []struct {
		name        string
		paramsJSON  string
		key         string
		returnError error
	}{
		{
			name:        "key not found",
			paramsJSON:  `{"key": "value"}`,
			key:         "test",
			returnError: config.ErrParamNotFound,
		},
		{
			name:        "not a string",
			paramsJSON:  `{"key": 1}`,
			key:         "key",
			returnError: config.ErrParamNotString,
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
			if err == nil {
				t.Fatalf("GetString expected to fail but succeeded, got: %v", value)
			}
			if !errors.Is(err, testCase.returnError) {
				t.Fatalf("GetString got error '%s', expected '%s'", err, testCase.returnError)
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

func TestBadIntParamsJSON(t *testing.T) {
	testCases := []struct {
		name        string
		paramsJSON  string
		key         string
		returnError error
	}{
		{
			name:        "key not found",
			paramsJSON:  `{"key": 1}`,
			key:         "test",
			returnError: config.ErrParamNotFound,
		},
		{
			name:        "not a number",
			paramsJSON:  `{"key": "1"}`,
			key:         "key",
			returnError: config.ErrParamNotNumber,
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
			if err == nil {
				t.Fatalf("GetInt expected to fail but succeeded, got: %v", value)
			}
			if !errors.Is(err, testCase.returnError) {
				t.Fatalf("GetInt got error '%s', expected '%s'", err, testCase.returnError)
			}
		})
	}
}

func TestGoodFloat32ParamsJSON(t *testing.T) {
	testCases := []struct {
		name       string
		paramsJSON string
		key        string
		expected   float32
	}{
		{
			name:       "no decimal param",
			paramsJSON: `{"key": 1}`,
			key:        "key",
			expected:   1,
		},
		{
			name:       "float param",
			paramsJSON: `{"key": 1.23}`,
			key:        "key",
			expected:   1.23,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			params := config.Params{}
			err := json.Unmarshal([]byte(testCase.paramsJSON), &params)
			if err != nil {
				t.Fatalf("Failed to unmarshal params JSON: %v", err)
			}
			value, err := params.GetFloat32(testCase.key)
			if err != nil {
				t.Fatalf("GetFloat32 returned error: %v", err)
			}
			if value != testCase.expected {
				t.Fatalf("GetFloat32 got %f, expected %f", value, testCase.expected)
			}
		})
	}
}

func TestBadFloat32ParamsJSON(t *testing.T) {
	testCases := []struct {
		name        string
		paramsJSON  string
		key         string
		returnError error
	}{
		{
			name:        "key not found",
			paramsJSON:  `{"key": 1}`,
			key:         "test",
			returnError: config.ErrParamNotFound,
		},
		{
			name:        "not a number",
			paramsJSON:  `{"key": "1"}`,
			key:         "key",
			returnError: config.ErrParamNotNumber,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			params := config.Params{}
			err := json.Unmarshal([]byte(testCase.paramsJSON), &params)
			if err != nil {
				t.Fatalf("Failed to unmarshal params JSON: %v", err)
			}
			value, err := params.GetFloat32(testCase.key)
			if err == nil {
				t.Fatalf("GetFloat32 expected to fail but succeeded, got: %v", value)
			}
			if !errors.Is(err, testCase.returnError) {
				t.Fatalf("GetFloat32 got error '%s', expected '%s'", err, testCase.returnError)
			}
		})
	}
}

func TestGoodFloat64ParamsJSON(t *testing.T) {
	testCases := []struct {
		name       string
		paramsJSON string
		key        string
		expected   float64
	}{
		{
			name:       "no decimal param",
			paramsJSON: `{"key": 1}`,
			key:        "key",
			expected:   1,
		},
		{
			name:       "float param",
			paramsJSON: `{"key": 1.23}`,
			key:        "key",
			expected:   1.23,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			params := config.Params{}
			err := json.Unmarshal([]byte(testCase.paramsJSON), &params)
			if err != nil {
				t.Fatalf("Failed to unmarshal params JSON: %v", err)
			}
			value, err := params.GetFloat64(testCase.key)
			if err != nil {
				t.Fatalf("GetFloat64 returned error: %v", err)
			}
			if value != testCase.expected {
				t.Fatalf("GetFloat64 got %f, expected %f", value, testCase.expected)
			}
		})
	}
}

func TestBadFloat64ParamsJSON(t *testing.T) {
	testCases := []struct {
		name        string
		paramsJSON  string
		key         string
		returnError error
	}{
		{
			name:        "key not found",
			paramsJSON:  `{"key": 1}`,
			key:         "test",
			returnError: config.ErrParamNotFound,
		},
		{
			name:        "not a number",
			paramsJSON:  `{"key": "1"}`,
			key:         "key",
			returnError: config.ErrParamNotNumber,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			params := config.Params{}
			err := json.Unmarshal([]byte(testCase.paramsJSON), &params)
			if err != nil {
				t.Fatalf("Failed to unmarshal params JSON: %v", err)
			}
			value, err := params.GetFloat64(testCase.key)
			if err == nil {
				t.Fatalf("GetFloat64 expected to fail but succeeded, got: %v", value)
			}
			if !errors.Is(err, testCase.returnError) {
				t.Fatalf("GetFloat64 got error '%s', expected '%s'", err, testCase.returnError)
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

func TestBadBoolParamsJSON(t *testing.T) {
	testCases := []struct {
		name        string
		paramsJSON  string
		key         string
		returnError error
	}{
		{
			name:        "key not found",
			paramsJSON:  `{"key": 1}`,
			key:         "test",
			returnError: config.ErrParamNotFound,
		},
		{
			name:        "not a bool",
			paramsJSON:  `{"key": "1"}`,
			key:         "key",
			returnError: config.ErrParamNotBool,
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
			if err == nil {
				t.Fatalf("GetBool expected to fail but succeeded, got: %v", value)
			}
			if !errors.Is(err, testCase.returnError) {
				t.Fatalf("GetBool got error '%s', expected '%s'", err, testCase.returnError)
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

func TestBadStringSliceParamsJSON(t *testing.T) {
	testCases := []struct {
		name        string
		paramsJSON  string
		key         string
		returnError error
	}{
		{
			name:        "key not found",
			paramsJSON:  `{"key": ["value1", "value2"]}`,
			key:         "test",
			returnError: config.ErrParamNotFound,
		},
		{
			name:        "not a slice",
			paramsJSON:  `{"key": "value"}`,
			key:         "key",
			returnError: config.ErrParamNotSlice,
		},
		{
			name:        "not a string slice",
			paramsJSON:  `{"key": ["value1", 2]}`,
			key:         "key",
			returnError: config.ErrParamNotStringSlice,
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
			if err == nil {
				t.Fatalf("GetStringSlice expected to fail but succeeded, got: %v", value)
			}
			if !errors.Is(err, testCase.returnError) {
				t.Fatalf("GetStringSlice got error '%s', expected '%s'", err, testCase.returnError)
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
		{
			name:       "int array with floats",
			paramsJSON: `{"key": [1.0, 2.0, 3.0]}`,
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

func TestBadIntSliceParamsJSON(t *testing.T) {
	testCases := []struct {
		name        string
		paramsJSON  string
		key         string
		returnError error
	}{
		{
			name:        "key not found",
			paramsJSON:  `{"key": ["value1", "value2"]}`,
			key:         "test",
			returnError: config.ErrParamNotFound,
		},
		{
			name:        "not a slice",
			paramsJSON:  `{"key": "value"}`,
			key:         "key",
			returnError: config.ErrParamNotIntSlice,
		},
		{
			name:        "not a int slice",
			paramsJSON:  `{"key": ["value1", 2]}`,
			key:         "key",
			returnError: config.ErrParamNotIntSlice,
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
			if err == nil {
				t.Fatalf("GetIntSlice expected to fail but succeeded, got: %v", value)
			}
			if !errors.Is(err, testCase.returnError) {
				t.Fatalf("GetIntSlice got error '%s', expected '%s'", err, testCase.returnError)
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

func TestBadByteSliceParamsJSON(t *testing.T) {
	testCases := []struct {
		name        string
		paramsJSON  string
		key         string
		returnError error
	}{
		{
			name:        "key not found",
			paramsJSON:  `{"key": ["value1", "value2"]}`,
			key:         "test",
			returnError: config.ErrParamNotFound,
		},
		{
			name:        "not a slice",
			paramsJSON:  `{"key": "value"}`,
			key:         "key",
			returnError: config.ErrParamNotByteSlice,
		},
		{
			name:        "not a int slice",
			paramsJSON:  `{"key": ["value1", 2]}`,
			key:         "key",
			returnError: config.ErrParamNotByteSlice,
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
			if err == nil {
				t.Fatalf("GetByteSlice expected to fail but succeeded, got: %v", value)
			}
			if !errors.Is(err, testCase.returnError) {
				t.Fatalf("GetByteSlice got error '%s', expected '%s'", err, testCase.returnError)
			}
		})
	}
}
