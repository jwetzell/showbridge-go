package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestFloatParseFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["float.parse"]
	if !ok {
		t.Fatalf("float.parse processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "float.parse",
	})

	if err != nil {
		t.Fatalf("failed to create float.parse processor: %s", err)
	}

	if processorInstance.Type() != "float.parse" {
		t.Fatalf("float.parse processor has wrong type: %s", processorInstance.Type())
	}
}

func TestFloatParseBadConfigBitSizeString(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["float.parse"]
	if !ok {
		t.Fatalf("float.parse processor not registered")
	}

	_, err := registration.New(config.ProcessorConfig{
		Type: "float.parse",
		Params: map[string]any{
			"bitSize": "64",
		},
	})

	if err == nil {
		t.Fatalf("float.parse should have returned an error for bad bitSize config")
	}
}

func TestFloatParseGoodConfig(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["float.parse"]
	if !ok {
		t.Fatalf("float.parse processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "float.parse",
		Params: map[string]any{
			"bitSize": 64.0,
		},
	})

	if err != nil {
		t.Fatalf("float.parse should have created processor but got error: %s", err)
	}

	payload := "12345.0"
	expected := float64(12345.0)

	got, err := processorInstance.Process(t.Context(), payload)
	if err != nil {
		t.Fatalf("float.parse processing failed: %s", err)
	}

	gotFloat, ok := got.(float64)
	if !ok {
		t.Fatalf("float.parse returned a %T payload: %s", got, got)
	}

	if gotFloat != expected {
		t.Fatalf("float.parse got %f, expected %f", gotFloat, expected)
	}
}

func TestGoodFloatParse(t *testing.T) {
	floatParser := processor.FloatParse{}
	tests := []struct {
		processor processor.Processor
		name      string
		payload   any
		bitSize   int
		expected  float64
	}{
		{
			name:     "positive number",
			payload:  "12345.67",
			bitSize:  64,
			expected: 12345.67,
		},
		{
			name:     "negative number",
			payload:  "-12345.67",
			bitSize:  64,
			expected: -12345.67,
		},
		{
			name:     "zero",
			payload:  "0",
			bitSize:  64,
			expected: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := floatParser.Process(t.Context(), test.payload)

			gotFloat, ok := got.(float64)
			if !ok {
				t.Fatalf("float.parse returned a %T payload: %s", got, got)
			}
			if err != nil {
				t.Fatalf("float.parse failed: %s", err)
			}
			if gotFloat != test.expected {
				t.Fatalf("float.parse got %f, expected %f", gotFloat, test.expected)
			}
		})
	}
}

func TestBadFloatParse(t *testing.T) {
	tests := []struct {
		processor   processor.Processor
		name        string
		payload     any
		bitSize     int
		errorString string
	}{
		{
			name:        "non-string input",
			payload:     []byte{0x01},
			bitSize:     64,
			errorString: "float.parse processor only accepts a string",
		},
		{
			name:        "not float string",
			payload:     "abcd",
			bitSize:     64,
			errorString: "strconv.ParseFloat: parsing \"abcd\": invalid syntax",
		},
		{
			name:        "bit size overflow",
			payload:     "1.79e+64",
			bitSize:     32,
			errorString: "strconv.ParseFloat: parsing \"1.79e+64\": value out of range",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			floatParser := processor.FloatParse{
				BitSize: test.bitSize,
			}
			got, err := floatParser.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("float.parse expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("float.parse got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
