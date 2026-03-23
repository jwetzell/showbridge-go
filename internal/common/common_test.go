package common_test

import (
	"slices"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
)

func TestGoodGetAnyAsInt(t *testing.T) {
	testCases := []struct {
		name       string
		value      any
		typedValue int
	}{
		{
			name:       "int",
			value:      int(42),
			typedValue: 42,
		},
		{
			name:       "uint",
			value:      uint(42),
			typedValue: 42,
		},
		{
			name:       "float32 without decimal",
			value:      float32(42.0),
			typedValue: 42,
		},
		{
			name:       "float64 without decimal",
			value:      float64(42.0),
			typedValue: 42,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value, ok := common.GetAnyAsInt(testCase.value)
			if !ok {
				t.Fatalf("GetAnyAsInt expected to succeed but failed")
			}
			if value != testCase.typedValue {
				t.Fatalf("GetAnyAsInt expected got %d,  expected %d", value, testCase.typedValue)
			}
		})
	}
}

func TestBadGetAnyAsInt(t *testing.T) {
	testCases := []struct {
		name  string
		value any
	}{
		{
			name:  "string",
			value: "value",
		},
		{
			name:  "float32 with decimal",
			value: float32(1.5),
		},
		{
			name:  "float64 with decimal",
			value: float64(1.5),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value, ok := common.GetAnyAsInt(testCase.value)
			if ok {
				t.Fatalf("GetAnyAsInt expected to fail but succeeded, got: %v", value)
			}
		})
	}
}

func TestGoodGetAnyAsByte(t *testing.T) {
	testCases := []struct {
		name       string
		value      any
		typedValue byte
	}{
		{
			name:       "int",
			value:      int(42),
			typedValue: 42,
		},
		{
			name:       "uint",
			value:      uint(42),
			typedValue: 42,
		},
		{
			name:       "float32 without decimal",
			value:      float32(42.0),
			typedValue: 42,
		},
		{
			name:       "float64 without decimal",
			value:      float64(42.0),
			typedValue: 42,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value, ok := common.GetAnyAsByte(testCase.value)
			if !ok {
				t.Fatalf("GetAnyAsByte expected to succeed but failed")
			}
			if value != testCase.typedValue {
				t.Fatalf("GetAnyAsByte expected got %d,  expected %d", value, testCase.typedValue)
			}
		})
	}
}

func TestBadGetAnyAsByte(t *testing.T) {
	testCases := []struct {
		name  string
		value any
	}{
		{
			name:  "string",
			value: "value",
		},
		{
			name:  "float32 with decimal",
			value: float32(1.5),
		},
		{
			name:  "float64 with decimal",
			value: float64(1.5),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value, ok := common.GetAnyAsByte(testCase.value)
			if ok {
				t.Fatalf("GetAnyAsByte expected to fail but succeeded, got: %v", value)
			}
		})
	}
}

func TestGoodGetAnyAsByteSlice(t *testing.T) {
	testCases := []struct {
		name       string
		value      any
		typedValue []byte
	}{
		{
			name:       "byte slice",
			value:      []byte{1, 2, 3},
			typedValue: []byte{1, 2, 3},
		},
		{
			name:       "int slice",
			value:      []int{1, 2, 3},
			typedValue: []byte{1, 2, 3},
		},
		{
			name:       "uint slice",
			value:      []uint{1, 2, 3},
			typedValue: []byte{1, 2, 3},
		},
		{
			name:       "float32 without decimal slice",
			value:      []float32{1, 2, 3},
			typedValue: []byte{1, 2, 3},
		},
		{
			name:       "float64 without decimal slice",
			value:      []float64{1, 2, 3},
			typedValue: []byte{1, 2, 3},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value, ok := common.GetAnyAsByteSlice(testCase.value)
			if !ok {
				t.Fatalf("GetAnyAsByteSlice expected to succeed but failed")
			}
			if !slices.Equal(value, testCase.typedValue) {
				t.Fatalf("GetAnyAsByteSlice expected got %d,  expected %d", value, testCase.typedValue)
			}
		})
	}
}

func TestBadGetAnyAsByteSlice(t *testing.T) {
	testCases := []struct {
		name  string
		value any
	}{
		{
			name:  "not a slice",
			value: "value",
		},
		{
			name:  "not a int slice",
			value: []any{"value1", 2},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value, ok := common.GetAnyAsByteSlice(testCase.value)
			if ok {
				t.Fatalf("GetAnyAsByteSlice expected to fail but succeeded, got: %v", value)
			}
		})
	}
}

func TestGoodGetAnyAsIntSlice(t *testing.T) {
	testCases := []struct {
		name       string
		value      any
		typedValue []int
	}{
		{
			name:       "int slice",
			value:      []int{1, 2, 3},
			typedValue: []int{1, 2, 3},
		},
		{
			name:       "byte slice",
			value:      []byte{1, 2, 3},
			typedValue: []int{1, 2, 3},
		},
		{
			name:       "uint slice",
			value:      []uint{1, 2, 3},
			typedValue: []int{1, 2, 3},
		},
		{
			name:       "float32 without decimal slice",
			value:      []float32{1, 2, 3},
			typedValue: []int{1, 2, 3},
		},
		{
			name:       "float64 without decimal slice",
			value:      []float64{1, 2, 3},
			typedValue: []int{1, 2, 3},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value, ok := common.GetAnyAsIntSlice(testCase.value)
			if !ok {
				t.Fatalf("GetAnyAsIntSlice expected to succeed but failed")
			}
			if !slices.Equal(value, testCase.typedValue) {
				t.Fatalf("GetAnyAsIntSlice expected got %d,  expected %d", value, testCase.typedValue)
			}
		})
	}
}

func TestBadGetAnyAsIntSlice(t *testing.T) {
	testCases := []struct {
		name  string
		value any
	}{
		{
			name:  "not a slice",
			value: "value",
		},
		{
			name:  "not a int slice",
			value: []any{"value1", 2},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value, ok := common.GetAnyAsIntSlice(testCase.value)
			if ok {
				t.Fatalf("GetAnyAsIntSlice expected to fail but succeeded, got: %v", value)
			}
		})
	}
}

func TestGoodGetAnyAsFloat32(t *testing.T) {
	testCases := []struct {
		name       string
		value      any
		typedValue float32
	}{
		{
			name:       "int",
			value:      int(42),
			typedValue: 42,
		},
		{
			name:       "uint",
			value:      uint(42),
			typedValue: 42,
		},
		{
			name:       "byte",
			value:      byte(42),
			typedValue: 42,
		},
		{
			name:       "float32",
			value:      float32(42.3),
			typedValue: 42.3,
		},
		{
			name:       "float64",
			value:      float64(42.3),
			typedValue: 42.3,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value, ok := common.GetAnyAsFloat32(testCase.value)
			if !ok {
				t.Fatalf("GetAnyAsFloat32 expected to succeed but failed")
			}
			if value != testCase.typedValue {
				t.Fatalf("GetAnyAsFloat32 expected got %f,  expected %f", value, testCase.typedValue)
			}
		})
	}
}

func TestBadGetAnyAsFloat32(t *testing.T) {
	testCases := []struct {
		name  string
		value any
	}{
		{
			name:  "string",
			value: "value",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value, ok := common.GetAnyAsFloat32(testCase.value)
			if ok {
				t.Fatalf("GetAnyAsFloat32 expected to fail but succeeded, got: %v", value)
			}
		})
	}
}

func TestGoodGetAnyAsFloat64(t *testing.T) {
	testCases := []struct {
		name       string
		value      any
		typedValue float64
	}{
		{
			name:       "int",
			value:      int(42),
			typedValue: 42,
		},
		{
			name:       "uint",
			value:      uint(42),
			typedValue: 42,
		},
		{
			name:       "byte",
			value:      byte(42),
			typedValue: 42,
		},
		{
			name:       "float32",
			value:      float32(42.5),
			typedValue: 42.5,
		},
		{
			name:       "float64",
			value:      float64(42.5),
			typedValue: 42.5,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value, ok := common.GetAnyAsFloat64(testCase.value)
			if !ok {
				t.Fatalf("GetAnyAsFloat64 expected to succeed but failed")
			}
			if value != testCase.typedValue {
				t.Fatalf("GetAnyAsFloat64 expected got %f,  expected %f", value, testCase.typedValue)
			}
		})
	}
}

func TestBadGetAnyAsFloat64(t *testing.T) {
	testCases := []struct {
		name  string
		value any
	}{
		{
			name:  "string",
			value: "value",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value, ok := common.GetAnyAsFloat64(testCase.value)
			if ok {
				t.Fatalf("GetAnyAsFloat64 expected to fail but succeeded, got: %v", value)
			}
		})
	}
}
