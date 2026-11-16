package types_test

import (
	"testing"

	"github.com/jwetzell/osc-go"
	"github.com/jwetzell/showbridge-go/internal/types"
	"github.com/stretchr/testify/require"
)

func TestGoodGetAnyAsByte(t *testing.T) {
	testCases := []struct {
		name     string
		input    any
		expected byte
	}{
		{
			name:     "byte",
			input:    byte('t'),
			expected: byte('t'),
		},
		{
			name:     "int",
			input:    int(1),
			expected: byte(1),
		},
		{
			name:     "uint",
			input:    uint(1),
			expected: byte(1),
		},
		{
			name:     "float32",
			input:    float32(1),
			expected: byte(1),
		},
		{
			name:     "float64",
			input:    float64(1),
			expected: byte(1),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, ok := types.GetAnyAsByte(tc.input)

			require.True(t, ok)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestBadGetAnyAsByte(t *testing.T) {
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
			_, ok := types.GetAnyAsByte(tc.input)

			require.False(t, ok)
		})
	}
}

func TestGoodGetAnyAsByteSlice(t *testing.T) {
	testCases := []struct {
		name     string
		input    any
		expected []byte
	}{
		{
			name:     "byte slice",
			input:    []byte{'t'},
			expected: []byte{'t'},
		},
		{
			name:     "string",
			input:    "test",
			expected: []byte("test"),
		},
		{
			name:     "int slice",
			input:    []int{1, 2, 3},
			expected: []byte{1, 2, 3},
		},
		{
			name:     "BinaryMarshaler",
			input:    osc.Message{Address: "/test", Args: []osc.Arg{}},
			expected: []byte{0x2f, 0x74, 0x65, 0x73, 0x74, 0x00, 0x00, 0x00, 0x2c, 0x00, 0x00, 0x00},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, ok := types.GetAnyAsByteSlice(tc.input)

			require.True(t, ok)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestBadGetAnyAsByteSlice(t *testing.T) {
	testCases := []struct {
		name  string
		input any
	}{
		{
			name:  "not a slice",
			input: int(1),
		},
		{
			name:  "BinaryMarshaler error",
			input: osc.Message{Address: "test", Args: []osc.Arg{}},
		},
		{
			name:  "not a slice of byte-able things",
			input: []string{"a", "b", "c"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, ok := types.GetAnyAsByteSlice(tc.input)

			require.False(t, ok)
		})
	}
}
