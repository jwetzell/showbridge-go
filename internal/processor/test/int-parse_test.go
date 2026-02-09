package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestIntParseFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["int.parse"]
	if !ok {
		t.Fatalf("int.parse processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "int.parse",
	})

	if err != nil {
		t.Fatalf("failed to create int.parse processor: %s", err)
	}

	if processorInstance.Type() != "int.parse" {
		t.Fatalf("int.parse processor has wrong type: %s", processorInstance.Type())
	}
}

func TestIntParseBadConfigBaseString(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["int.parse"]
	if !ok {
		t.Fatalf("int.parse processor not registered")
	}

	_, err := registration.New(config.ProcessorConfig{
		Type: "int.parse",
		Params: map[string]any{
			"base": "10",
		},
	})

	if err == nil {
		t.Fatalf("int.parse should have returned an error for bad base config")
	}
}

func TestIntParseBadConfigBitSizeString(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["int.parse"]
	if !ok {
		t.Fatalf("int.parse processor not registered")
	}

	_, err := registration.New(config.ProcessorConfig{
		Type: "int.parse",
		Params: map[string]any{
			"bitSize": "64",
		},
	})

	if err == nil {
		t.Fatalf("int.parse should have returned an error for bad bitSize config")
	}
}

func TestIntParseGoodConfig(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["int.parse"]
	if !ok {
		t.Fatalf("int.parse processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "int.parse",
		Params: map[string]any{
			"base":    10.0,
			"bitSize": 64.0,
		},
	})

	if err != nil {
		t.Fatalf("int.parse should have created processor but got error: %s", err)
	}

	payload := "12345"
	expected := int64(12345)

	got, err := processorInstance.Process(t.Context(), payload)
	if err != nil {
		t.Fatalf("int.parse processing failed: %s", err)
	}

	gotInt, ok := got.(int64)
	if !ok {
		t.Fatalf("int.parse returned a %T payload: %s", got, got)
	}

	if gotInt != expected {
		t.Fatalf("int.parse got %d, expected %d", gotInt, expected)
	}
}

func TestGoodIntParse(t *testing.T) {
	tests := []struct {
		processor processor.Processor
		name      string
		payload   any
		expected  int64
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
			name:     "negative number",
			payload:  "-12345",
			expected: -12345,
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
			intParser := processor.IntParse{
				Base:    test.base,
				BitSize: test.bitSize,
			}

			got, err := intParser.Process(t.Context(), test.payload)

			gotInt, ok := got.(int64)
			if !ok {
				t.Fatalf("int.parse returned a %T payload: %s", got, got)
			}
			if err != nil {
				t.Fatalf("int.parse failed: %s", err)
			}
			if gotInt != test.expected {
				t.Fatalf("int.parse got %d, expected %d", gotInt, test.expected)
			}
		})
	}
}

func TestBadIntParse(t *testing.T) {
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
			errorString: "int.parse processor only accepts a string",
		},
		{
			name:        "not int string",
			payload:     "123.46",
			base:        10,
			bitSize:     64,
			errorString: "strconv.ParseInt: parsing \"123.46\": invalid syntax",
		},
		{
			name:        "bit overflow",
			payload:     "12345678901234567890",
			base:        10,
			bitSize:     32,
			errorString: "strconv.ParseInt: parsing \"12345678901234567890\": value out of range",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			intParser := processor.IntParse{
				Base:    test.base,
				BitSize: test.bitSize,
			}
			got, err := intParser.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("int.parse expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("int.parse got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
