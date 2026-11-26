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

type TestMIDIOutputModule struct {
	outputFunc func(ctx context.Context, msg midi.Message) error
}

func (m *TestMIDIOutputModule) Start(ctx context.Context, router routing.Router) error {
	<-ctx.Done()
	return nil
}

func (m *TestMIDIOutputModule) Stop() {}

func (m *TestMIDIOutputModule) Output(ctx context.Context, msg midi.Message) error {
	if m.outputFunc != nil {
		return m.outputFunc(ctx, msg)
	}
	return nil
}

func NewTestMIDIOutputModule(outputFunc func(ctx context.Context, msg midi.Message) error) *TestMIDIOutputModule {
	return &TestMIDIOutputModule{
		outputFunc: outputFunc,
	}
}

func TestGoodMIDIControlChange(t *testing.T) {
	testCases := []struct {
		name     string
		params   map[string]any
		payload  any
		expected midi.Message
	}{
		{
			name: "simple control change",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"control": "64",
				"value":   "65",
			},
			payload:  nil,
			expected: midi.ControlChange(1, 64, 65),
		},
		{
			name: "templated channel",
			params: map[string]any{
				"module":  "test",
				"channel": "{{.Payload}}",
				"control": "64",
				"value":   "65",
			},
			payload:  1,
			expected: midi.ControlChange(1, 64, 65),
		},
		{
			name: "templated control",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"control": "{{.Payload}}",
				"value":   "65",
			},
			payload:  64,
			expected: midi.ControlChange(1, 64, 65),
		},
		{
			name: "templated value",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"control": "64",
				"value":   "{{.Payload}}",
			},
			payload:  65,
			expected: midi.ControlChange(1, 64, 65),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("midi.control_change")
			require.True(t, ok, "processor registration not found")

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "midi.control_change",
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

func TestBadMIDIControlChange(t *testing.T) {
	testCases := []struct {
		name          string
		params        map[string]any
		payload       any
		moduleManager routing.ModuleManager
		errorString   string
	}{
		{
			name: "control change no module",
			params: map[string]any{
				"channel": "1",
				"control": "64",
				"value":   "64",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "midi.control_change module error: not found",
		},
		{
			name: "control change module not a string",
			params: map[string]any{
				"module":  123,
				"channel": "1",
				"control": "64",
				"value":   "64",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "midi.control_change module error: not a string",
		},
		{
			name: "control change no channel",
			params: map[string]any{
				"module":  "test",
				"control": "64",
				"value":   "64",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "midi.control_change channel error: not found",
		},
		{
			name: "control change channel not a string",
			params: map[string]any{
				"module":  "test",
				"channel": 123,
				"control": "64",
				"value":   "64",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "midi.control_change channel error: not a string",
		},
		{
			name: "control change no control",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"value":   "64",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "midi.control_change control error: not found",
		},
		{
			name: "control change control not a string",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"control": 123,
				"value":   "64",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "midi.control_change control error: not a string",
		},
		{
			name: "control change no value",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"control": "64",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "midi.control_change value error: not found",
		},
		{
			name: "control change value not a string",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"control": "64",
				"value":   123,
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "midi.control_change value error: not a string",
		},
		{
			name: "channel template syntax error",
			params: map[string]any{
				"module":  "test",
				"channel": "{{.Payload",
				"control": "64",
				"value":   "127",
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
			name: "control template syntax error",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"control": "{{.Payload",
				"value":   "127",
			},
			payload: 1,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "template: control:1: unclosed action",
		},
		{
			name: "value template syntax error",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"control": "64",
				"value":   "{{.Payload",
			},
			payload: 1,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "template: value:1: unclosed action",
		},
		{
			name: "channel template undefined variable",
			params: map[string]any{
				"module":  "test",
				"channel": "{{.Undefined}}",
				"control": "64",
				"value":   "127",
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
			name: "control template undefined variable",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"control": "{{.Undefined}}",
				"value":   "127",
			},
			payload: 1,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "template: control:1:2: executing \"control\" at <.Undefined>: can't evaluate field Undefined in type processor.templateData",
		},
		{
			name: "value template undefined variable",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"control": "64",
				"value":   "{{.Undefined}}",
			},
			payload: 1,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "template: value:1:2: executing \"value\" at <.Undefined>: can't evaluate field Undefined in type processor.templateData",
		},
		{
			name: "channel not parseable",
			params: map[string]any{
				"module":  "test",
				"channel": "abc",
				"control": "64",
				"value":   "127",
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
			name: "control not parseable",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"control": "abc",
				"value":   "127",
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
			name: "value not parseable",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"control": "64",
				"value":   "abc",
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
				"control": "64",
				"value":   "127",
			},
			payload:       nil,
			moduleManager: nil,
			errorString:   "midi.control_change module manager is nil",
		},
		{
			name: "missing module",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"control": "64",
				"value":   "127",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"asdf": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					require.Fail(t, "unexpected call to MIDI output module")
					return nil
				}),
			},
			errorString: "midi.control_change unable to find module with id: test",
		},
		{
			name: "module not a MIDIOutputModule",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"control": "64",
				"value":   "127",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestPubSubModule("test", nil),
			},
			errorString: "midi.control_change module with id test is not a MIDIOutputModule",
		},
		{
			name: "module output errored",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"control": "64",
				"value":   "127",
			},
			payload: nil,
			moduleManager: testModuleManager{
				"test": NewTestMIDIOutputModule(func(ctx context.Context, msg midi.Message) error {
					return errors.New("output error")
				}),
			},
			errorString: "midi.control_change failed to send output: output error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("midi.control_change")
			require.True(t, ok, "processor registration not found")

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "midi.control_change",
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
