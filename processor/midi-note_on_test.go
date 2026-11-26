package processor_test

import (
	"context"
	"errors"
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
		{
			name: "templated channel",
			params: map[string]any{
				"module":   "test",
				"channel":  "{{.Payload}}",
				"note":     "64",
				"velocity": "127",
			},
			payload:  1,
			expected: midi.NoteOn(1, 64, 127),
		},
		{
			name: "templated note",
			params: map[string]any{
				"module":   "test",
				"channel":  "1",
				"note":     "{{.Payload}}",
				"velocity": "127",
			},
			payload:  64,
			expected: midi.NoteOn(1, 64, 127),
		},
		{
			name: "templated velocity",
			params: map[string]any{
				"module":   "test",
				"channel":  "1",
				"note":     "64",
				"velocity": "{{.Payload}}",
			},
			payload:  127,
			expected: midi.NoteOn(1, 64, 127),
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
		name          string
		params        map[string]any
		payload       any
		moduleManager routing.ModuleManager
		errorString   string
	}{
		{
			name: "note off no module",
			params: map[string]any{
				"channel":  "1",
				"note":     "64",
				"velocity": "64",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
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
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "midi.note_on module error: not a string",
		},
		{
			name: "note off no channel",
			params: map[string]any{
				"module":   "test",
				"note":     "64",
				"velocity": "64",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
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
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "midi.note_on channel error: not a string",
		},
		{
			name: "note off no note",
			params: map[string]any{
				"module":   "test",
				"channel":  "1",
				"velocity": "64",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
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
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "midi.note_on note error: not a string",
		},
		{
			name: "note off no velocity",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"note":    "64",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
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
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "midi.note_on velocity error: not a string",
		},
		{
			name: "channel template syntax error",
			params: map[string]any{
				"module":   "test",
				"channel":  "{{.Payload",
				"note":     "64",
				"velocity": "127",
			},
			payload: 1,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "template: channel:1: unclosed action",
		},
		{
			name: "note template syntax error",
			params: map[string]any{
				"module":   "test",
				"channel":  "1",
				"note":     "{{.Payload",
				"velocity": "127",
			},
			payload: 1,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "template: note:1: unclosed action",
		},
		{
			name: "velocity template syntax error",
			params: map[string]any{
				"module":   "test",
				"channel":  "1",
				"note":     "64",
				"velocity": "{{.Payload",
			},
			payload: 1,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "template: velocity:1: unclosed action",
		},
		{
			name: "channel template undefined variable",
			params: map[string]any{
				"module":   "test",
				"channel":  "{{.Undefined}}",
				"note":     "64",
				"velocity": "127",
			},
			payload: 1,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "template: channel:1:2: executing \"channel\" at <.Undefined>: can't evaluate field Undefined in type processor.templateData",
		},
		{
			name: "note template undefined variable",
			params: map[string]any{
				"module":   "test",
				"channel":  "1",
				"note":     "{{.Undefined}}",
				"velocity": "127",
			},
			payload: 1,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "template: note:1:2: executing \"note\" at <.Undefined>: can't evaluate field Undefined in type processor.templateData",
		},
		{
			name: "velocity template undefined variable",
			params: map[string]any{
				"module":   "test",
				"channel":  "1",
				"note":     "64",
				"velocity": "{{.Undefined}}",
			},
			payload: 1,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "template: velocity:1:2: executing \"velocity\" at <.Undefined>: can't evaluate field Undefined in type processor.templateData",
		},
		{
			name: "channel not parseable",
			params: map[string]any{
				"module":   "test",
				"channel":  "abc",
				"note":     "64",
				"velocity": "127",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "strconv.ParseUint: parsing \"abc\": invalid syntax",
		},
		{
			name: "note not parseable",
			params: map[string]any{
				"module":   "test",
				"channel":  "1",
				"note":     "abc",
				"velocity": "127",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "strconv.ParseUint: parsing \"abc\": invalid syntax",
		},
		{
			name: "velocity not parseable",
			params: map[string]any{
				"module":   "test",
				"channel":  "1",
				"note":     "64",
				"velocity": "abc",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "strconv.ParseUint: parsing \"abc\": invalid syntax",
		},
		{
			name: "missing module manager",
			params: map[string]any{
				"module":   "test",
				"channel":  "1",
				"note":     "64",
				"velocity": "127",
			},
			payload:       nil,
			moduleManager: nil,
			errorString:   "midi.note_on module manager is nil",
		},
		{
			name: "missing module",
			params: map[string]any{
				"module":   "test",
				"channel":  "1",
				"note":     "64",
				"velocity": "127",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"asdf": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "midi.note_on unable to find module with id: test",
		},
		{
			name: "module not a MIDIOutputModule",
			params: map[string]any{
				"module":   "test",
				"channel":  "1",
				"note":     "64",
				"velocity": "127",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestPubSubModule("test", nil),
			},
			errorString: "midi.note_on module with id test is not a MIDIOutputModule",
		},
		{
			name: "module output errored",
			params: map[string]any{
				"module":   "test",
				"channel":  "1",
				"note":     "64",
				"velocity": "127",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					return errors.New("output error")
				}),
			},
			errorString: "midi.note_on failed to send output: output error",
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

			_, err = EvaluateProcessFunc(processorInstance, t.Context(), routing.WrappedPayload{Payload: tc.payload}, tc.moduleManager, nil)
			require.EqualError(t, err, tc.errorString)
		})
	}
}
