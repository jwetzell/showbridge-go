package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

func TestTimeSleepFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["time.sleep"]
	if !ok {
		t.Fatalf("time.sleep processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "time.sleep",
		Params: map[string]any{
			"duration": 1000.0,
		},
	})

	if err != nil {
		t.Fatalf("failed to create time.sleep processor: %s", err)
	}

	if processorInstance.Type() != "time.sleep" {
		t.Fatalf("time.sleep processor has wrong type: %s", processorInstance.Type())
	}
}

func TestGoodTimeSleep(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]any
		payload any
	}{
		{
			name:    "string payload",
			payload: "hello",
			params:  map[string]any{"duration": 100.0},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["time.sleep"]
			if !ok {
				t.Fatalf("time.sleep processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "time.sleep",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("time.sleep failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err != nil {
				t.Fatalf("time.sleep failed: %s", err)
			}

			if got != test.payload {
				t.Fatalf("time.sleep got %+v, expected %+v", got, test.payload)
			}
		})
	}
}

func TestBadTimeSleep(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name:        "no-duration param",
			payload:     "hello",
			params:      map[string]any{},
			errorString: "time.sleep requires a duration parameter",
		},
		{
			name:    "non-number duration param",
			payload: "hello",
			params: map[string]any{
				"duration": "1000",
			},
			errorString: "time.sleep duration must be a number",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["time.sleep"]
			if !ok {
				t.Fatalf("time.sleep processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "time.sleep",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("string.split got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), test.payload)

			if err == nil {
				t.Fatalf("time.sleep expected to fail but succeeded, got: %v", got)

			}
			if err.Error() != test.errorString {
				t.Fatalf("time.sleep got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
