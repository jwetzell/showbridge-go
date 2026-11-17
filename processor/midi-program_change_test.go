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
		name          string
		params        map[string]any
		payload       any
		moduleManager routing.ModuleManager
		errorString   string
	}{
		{
			name: "program change no module",
			params: map[string]any{
				"channel": "1",
				"program": "64",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "midi.program_change module error: not found",
		},
		{
			name: "program change module not a string",
			params: map[string]any{
				"module":  123,
				"channel": "1",
				"program": "64",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "midi.program_change module error: not a string",
		},
		{
			name: "program change no channel",
			params: map[string]any{
				"module":  "test",
				"program": "64",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "midi.program_change channel error: not found",
		},
		{
			name: "program change channel not a string",
			params: map[string]any{
				"module":  "test",
				"channel": 123,
				"program": "64",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "midi.program_change channel error: not a string",
		},
		{
			name: "program change no program",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "midi.program_change program error: not found",
		},
		{
			name: "program change program not a string",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"program": 123,
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "midi.program_change program error: not a string",
		},
		{
			name: "channel template syntax error",
			params: map[string]any{
				"module":  "test",
				"channel": "{{.Payload",
				"program": "64",
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
			name: "program template syntax error",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"program": "{{.Payload",
			},
			payload: 1,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "template: program:1: unclosed action",
		},
		{
			name: "channel template undefined variable",
			params: map[string]any{
				"module":  "test",
				"channel": "{{.Undefined}}",
				"program": "64",
			},
			payload: 1,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "template: channel:1:2: executing \"channel\" at <.Undefined>: can't evaluate field Undefined in type routing.WrappedPayload",
		},
		{
			name: "program template undefined variable",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"program": "{{.Undefined}}",
			},
			payload: 1,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "template: program:1:2: executing \"program\" at <.Undefined>: can't evaluate field Undefined in type routing.WrappedPayload",
		},
		{
			name: "channel not parseable",
			params: map[string]any{
				"module":  "test",
				"channel": "abc",
				"program": "64",
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
			name: "program not parseable",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"program": "abc",
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
				"module":  "test",
				"channel": "1",
				"program": "64",
			},
			payload:       nil,
			moduleManager: nil,
			errorString:   "midi.program_change module manager is nil",
		},
		{
			name: "missing module",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"program": "64",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"asdf": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "midi.program_change unable to find module with id: test",
		},
		{
			name: "module not a MIDIOutputModule",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"program": "64",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestPubSubModule("test"),
			},
			errorString: "midi.program_change module with id test is not a MIDIOutputModule",
		},
		{
			name: "module output errored",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"program": "64",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					return errors.New("output error")
				}),
			},
			errorString: "midi.program_change failed to send output: output error",
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

			_, err = EvaluateProcessFunc(processorInstance, t.Context(), routing.WrappedPayload{Payload: tc.payload}, tc.moduleManager, nil)
			require.EqualError(t, err, tc.errorString)
		})
	}
}
