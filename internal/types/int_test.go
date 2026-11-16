package types_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/types"
	"github.com/stretchr/testify/require"
)

func TestGoodGetAnyAsInt(t *testing.T) {
	testCases := []struct {
		name     string
		input    any
		expected int
	}{
		{
			name:     "byte",
			input:    byte('t'),
			expected: int('t'),
		},
		{
			name:     "int",
			input:    int(1),
			expected: int(1),
		},
		{
			name:     "uint",
			input:    uint(1),
			expected: int(1),
		},
		{
			name:     "float32",
			input:    float32(1),
			expected: int(1),
		},
		{
			name:     "float64",
			input:    float64(1),
			expected: int(1),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, ok := types.GetAnyAsInt(tc.input)

			require.True(t, ok)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestBadGetAnyAsInt(t *testing.T) {
	testCases := []struct {
		name  string
		input any
	}{
		{
			name:  "string",
			input: "test",
		},
		{
			name:  "float32 with decimal",
			input: float32(1.1),
		},
		{
			name:  "float64 with decimal",
			input: float64(1.1),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, ok := types.GetAnyAsInt(tc.input)

			require.False(t, ok)
		})
	}
}

func TestGoodGetAnyAsIntSlice(t *testing.T) {
	testCases := []struct {
		name     string
		input    any
		expected []int
	}{
		{
			name:     "byte slice",
			input:    []byte{0x00, 0x01, 0x02},
			expected: []int{0x00, 0x01, 0x02},
		},
		{
			name:     "int slice",
			input:    []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, ok := types.GetAnyAsIntSlice(tc.input)

			require.True(t, ok)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestBadGetAnyAsIntSlice(t *testing.T) {
	testCases := []struct {
		name  string
		input any
	}{
		{
			name:  "not a slice",
			input: int(1),
		},
		{
			name:  "not a slice of int-able things",
			input: []string{"a", "b", "c"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, ok := types.GetAnyAsIntSlice(tc.input)

			require.False(t, ok)
		})
	}
}
