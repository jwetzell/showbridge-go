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
		name        string
		params      map[string]any
		payload     any
		errorString string
	}{
		{
			name: "control change no module",
			params: map[string]any{
				"channel": "1",
				"control": "64",
				"value":   "64",
			},
			payload:     nil,
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
			payload:     nil,
			errorString: "midi.control_change module error: not a string",
		},
		{
			name: "control change no channel",
			params: map[string]any{
				"module":  "test",
				"control": "64",
				"value":   "64",
			},
			payload:     nil,
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
			payload:     nil,
			errorString: "midi.control_change channel error: not a string",
		},
		{
			name: "control change no control",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"value":   "64",
			},
			payload:     nil,
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
			payload:     nil,
			errorString: "midi.control_change control error: not a string",
		},
		{
			name: "control change no value",
			params: map[string]any{
				"module":  "test",
				"channel": "1",
				"control": "64",
			},
			payload:     nil,
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
			payload:     nil,
			errorString: "midi.control_change value error: not a string",
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
