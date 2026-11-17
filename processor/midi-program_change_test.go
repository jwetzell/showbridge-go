package processor_test

import (
	"context"
	"testing"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/processor"
	"github.com/jwetzell/showbridge-go/routing"
	"github.com/stretchr/testify/require"
	"gitlab.com/gomidi/midi/v2"
)

func TestGoodMIDIProgramChange(t *testing.T) {
	testCases := []struct {
		name     string
		params   map[string]any
		payload  any
		expected midi.Message
	}{
		{
			name: "simple program change",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"program": "64",
			},
			payload:  nil,
			expected: midi.ProgramChange(1, 64),
		},
		{
			name: "templated channel",
			params: map[string]any{
				"module":  "test",
				"channel": "{{.Payload}}",
				"program": "64",
			},
			payload:  1,
			expected: midi.ProgramChange(1, 64),
		},
		{
			name: "templated program",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"program": "{{.Payload}}",
			},
			payload:  64,
			expected: midi.ProgramChange(1, 64),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("midi.program_change")
			require.True(t, ok, "processor registration not found")

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "midi.program_change",
				Params: tc.params,
			})
			require.NoError(t, err, "failed to create processor instance")

			_, err = EvaluateProcessFunc(processorInstance, t.Context(), routing.WrappedPayload{Payload: tc.payload}, testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Equal(t, tc.expected, msg)
					return nil
				}),
			}, nil)
			require.NoError(t, err, "failed to evaluate process function")

		})
	}
}

func TestBadMIDIProgramChange(t *testing.T) {
	testCases := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name: "program change no module",
			params: map[string]any{
				"channel": "1",
				"program": "64",
			},
			payload:     nil,
			errorString: "midi.program_change module error: not found",
		},
		{
			name: "program change module not a string",
			params: map[string]any{
				"module":  123,
				"channel": "1",
				"program": "64",
			},
			payload:     nil,
			errorString: "midi.program_change module error: not a string",
		},
		{
			name: "program change no channel",
			params: map[string]any{
				"module":  "test",
				"program": "64",
			},
			payload:     nil,
			errorString: "midi.program_change channel error: not found",
		},
		{
			name: "program change channel not a string",
			params: map[string]any{
				"module":  "test",
				"channel": 123,
				"program": "64",
			},
			payload:     nil,
			errorString: "midi.program_change channel error: not a string",
		},
		{
			name: "program change no program",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
			},
			payload:     nil,
			errorString: "midi.program_change program error: not found",
		},
		{
			name: "program change program not a string",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"program": 123,
			},
			payload:     nil,
			errorString: "midi.program_change program error: not a string",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("midi.program_change")
			require.True(t, ok, "processor registration not found")

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "midi.program_change",
				Params: tc.params,
			})
			if err != nil {
				require.EqualError(t, err, tc.errorString)
				return
			}

			_, err = EvaluateProcessFunc(processorInstance, t.Context(), routing.WrappedPayload{Payload: tc.payload}, testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			}, nil)
			require.EqualError(t, err, tc.errorString)
		})
	}
}
