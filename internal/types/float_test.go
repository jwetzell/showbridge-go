package types_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/types"
	"github.com/stretchr/testify/require"
)

func TestGoodGetAnyAsFloat32(t *testing.T) {
	testCases := []struct {
		name     string
		input    any
		expected float32
	}{
		{
			name:     "byte",
			input:    byte(1),
			expected: float32(1),
		},
		{
			name:     "int",
			input:    int(1),
			expected: float32(1),
		},
		{
			name:     "uint",
			input:    uint(1),
			expected: float32(1),
		},
		{
			name:     "float32",
			input:    float32(1),
			expected: float32(1),
		},
		{
			name:     "float64",
			input:    float64(1),
			expected: float32(1),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, ok := types.GetAnyAsFloat32(tc.input)

			require.True(t, ok)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestBadGetAnyAsFloat32(t *testing.T) {
	testCases := []struct {
		name  string
		input any
	}{
		{
			name:  "string",
			input: "test",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, ok := types.GetAnyAsFloat32(tc.input)

			require.False(t, ok)
		})
	}
}

func TestGoodGetAnyAsFloat64(t *testing.T) {
	testCases := []struct {
		name     string
		input    any
		expected float64
	}{
		{
			name:     "byte",
			input:    byte(1),
			expected: float64(1),
		},
		{
			name:     "int",
			input:    int(1),
			expected: float64(1),
		},
		{
			name:     "uint",
			input:    uint(1),
			expected: float64(1),
		},
		{
			name:     "float32",
			input:    float32(1),
			expected: float64(1),
		},
		{
			name:     "float64",
			input:    float64(1),
			expected: float64(1),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, ok := types.GetAnyAsFloat64(tc.input)

			require.True(t, ok)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestBadGetAnyAsFloat64(t *testing.T) {
	testCases := []struct {
		name  string
		input any
	}{
		{
			name:  "string",
			input: "test",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, ok := types.GetAnyAsFloat64(tc.input)

			require.False(t, ok)
		})
	}
}
