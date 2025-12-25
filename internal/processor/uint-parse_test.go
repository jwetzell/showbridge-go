package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestUintParseFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["uint.parse"]
	if !ok {
		t.Fatalf("uint.parse processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "uint.parse",
	})

	if err != nil {
		t.Fatalf("failed to create uint.parse processor: %s", err)
	}

	if processorInstance.Type() != "uint.parse" {
		t.Fatalf("uint.parse processor has wrong type: %s", processorInstance.Type())
	}
}

func TestUintParseBadConfigBaseString(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["uint.parse"]
	if !ok {
		t.Fatalf("uint.parse processor not registered")
	}

	_, err := registration.New(config.ProcessorConfig{
		Type: "uint.parse",
		Params: map[string]any{
			"base": "10",
		},
	})

	if err == nil {
		t.Fatalf("uint.parse should have returned an error for bad base config")
	}
}

func TestUintParseBadConfigBitSizeString(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["uint.parse"]
	if !ok {
		t.Fatalf("uint.parse processor not registered")
	}

	_, err := registration.New(config.ProcessorConfig{
		Type: "uint.parse",
		Params: map[string]any{
			"bitSize": "64",
		},
	})

	if err == nil {
		t.Fatalf("uint.parse should have returned an error for bad bitSize config")
	}
}

func TestUintParseGoodConfig(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["uint.parse"]
	if !ok {
		t.Fatalf("uint.parse processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "uint.parse",
		Params: map[string]any{
			"base":    10.0,
			"bitSize": 64.0,
		},
	})

	if err != nil {
		t.Fatalf("uint.parse should have created processor but got error: %s", err)
	}

	payload := "12345"
	expected := uint64(12345)

	got, err := processorInstance.Process(t.Context(), payload)
	if err != nil {
		t.Fatalf("uint.parse processing failed: %s", err)
	}

	gotUint, ok := got.(uint64)
	if !ok {
		t.Fatalf("uint.parse returned a %T payload: %s", got, got)
	}

	if gotUint != expected {
		t.Fatalf("uint.parse got %d, expected %d", gotUint, expected)
	}
}

func TestGoodUintParse(t *testing.T) {

	tests := []struct {
		processor processor.Processor
		name      string
		payload   any
		expected  uint64
		base      int
		bitSize   int
	}{
		{
			name:     "positive number",
			payload:  "12345",
			expected: 12345,
			base:     10,
			bitSize:  64,
		},
		{
			name:     "zero",
			payload:  "0",
			expected: 0,
			base:     10,
			bitSize:  64,
		},
		{
			name:     "binary",
			payload:  "1010101",
			expected: 85,
			base:     2,
			bitSize:  64,
		},
		{
			name:     "hex",
			payload:  "15F",
			expected: 351,
			base:     16,
			bitSize:  64,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uintParser := processor.UintParse{
				Base:    test.base,
				BitSize: test.bitSize,
			}
			got, err := uintParser.Process(t.Context(), test.payload)

			gotUint, ok := got.(uint64)
			if !ok {
				t.Fatalf("uint.parse returned a %T payload: %s", got, got)
			}
			if err != nil {
				t.Fatalf("uint.parse failed: %s", err)
			}
			if gotUint != test.expected {
				t.Fatalf("uint.parse got %d, expected %d", gotUint, test.expected)
			}
		})
	}
}

func TestBadUintParse(t *testing.T) {
	uintParser := processor.UintParse{}
	tests := []struct {
		processor   processor.Processor
		name        string
		payload     any
		base        int
		bitSize     int
		errorString string
	}{
		{
			name:        "non-string input",
			payload:     []byte{0x01},
			base:        10,
			bitSize:     64,
			errorString: "uint.parse processor only accepts a string",
		},
		{
			name:        "not uint string",
			payload:     "-1234",
			base:        10,
			bitSize:     64,
			errorString: "strconv.ParseUint: parsing \"-1234\": invalid syntax",
		},
		{
			name:        "bit overflow",
			payload:     "123456789012345678901234567",
			base:        10,
			bitSize:     32,
			errorString: "strconv.ParseUint: parsing \"123456789012345678901234567\": value out of range",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := uintParser.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("uint.parse expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("uint.parse got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
