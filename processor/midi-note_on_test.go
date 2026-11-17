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

func TestGoodMIDINoteOn(t *testing.T) {
	testCases := []struct {
		name     string
		params   map[string]any
		payload  any
		expected midi.Message
	}{
		{
			name: "simple note on",
			params: map[string]any{
				"module":   "test",
				"channel":  "1",
				"note":     "64",
				"velocity": "64",
			},
			payload:  nil,
			expected: midi.NoteOn(1, 64, 64),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("midi.note_on")
			require.True(t, ok, "processor registration not found")

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "midi.note_on",
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

func TestBadMIDINoteOn(t *testing.T) {
	testCases := []struct {
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name: "note off no module",
			params: map[string]any{
				"channel":  "1",
				"note":     "64",
				"velocity": "64",
			},
			payload:     nil,
			errorString: "midi.note_on module error: not found",
		},
		{
			name: "note off module not a string",
			params: map[string]any{
				"module":   123,
				"channel":  "1",
				"note":     "64",
				"velocity": "64",
			},
			payload:     nil,
			errorString: "midi.note_on module error: not a string",
		},
		{
			name: "note off no channel",
			params: map[string]any{
				"module":   "test",
				"note":     "64",
				"velocity": "64",
			},
			payload:     nil,
			errorString: "midi.note_on channel error: not found",
		},
		{
			name: "note off channel not a string",
			params: map[string]any{
				"module":   "test",
				"channel":  123,
				"note":     "64",
				"velocity": "64",
			},
			payload:     nil,
			errorString: "midi.note_on channel error: not a string",
		},
		{
			name: "note off no note",
			params: map[string]any{
				"module":   "test",
				"channel":  "1",
				"velocity": "64",
			},
			payload:     nil,
			errorString: "midi.note_on note error: not found",
		},
		{
			name: "note off note not a string",
			params: map[string]any{
				"module":   "test",
				"channel":  "1",
				"note":     123,
				"velocity": "64",
			},
			payload:     nil,
			errorString: "midi.note_on note error: not a string",
		},
		{
			name: "note off no velocity",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"note":    "64",
			},
			payload:     nil,
			errorString: "midi.note_on velocity error: not found",
		},
		{
			name: "note off velocity not a string",
			params: map[string]any{
				"module":   "test",
				"channel":  "1",
				"note":     "64",
				"velocity": 123,
			},
			payload:     nil,
			errorString: "midi.note_on velocity error: not a string",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("midi.note_on")
			require.True(t, ok, "processor registration not found")

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "midi.note_on",
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
