package framer_test

import (
	"slices"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/framer"
)

func TestGoodRawFramerDecode(t *testing.T) {
	tests := []struct {
		name     string
		framer   framer.Framer
		input    []byte
		expected [][]byte
	}{
		{
			name:   "basic raw framer",
			framer: framer.GetFramer("RAW"),
			input:  []byte("Hello\nWorld\nThis is a test\n"),
			expected: [][]byte{
				[]byte("Hello\nWorld\nThis is a test\n"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			frames := test.framer.Decode(test.input)
			if len(frames) != len(test.expected) {
				t.Errorf("raw framer got %d frames, expected %d", len(frames), len(test.expected))
			}
			for i, frame := range frames {
				if !slices.Equal(frame, test.expected[i]) {
					t.Errorf("raw framer frame %d got %s, expected %s", i, frame, test.expected[i])
				}
			}
		})
	}
}

func TestGoodRawFramerEncode(t *testing.T) {
	tests := []struct {
		name     string
		framer   framer.Framer
		expected []byte
		input    []byte
	}{
		{
			name:     "basic raw framer",
			framer:   framer.GetFramer("RAW"),
			expected: []byte("Hello\nWorld\nThis is a test\n"),
			input:    []byte("Hello\nWorld\nThis is a test\n"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			frame := test.framer.Encode(test.input)
			if len(frame) != len(test.expected) {
				t.Errorf("raw framer got %d frames, expected %d", len(frame), len(test.expected))
			}
			if !slices.Equal(frame, test.expected) {
				t.Errorf("raw frame got %s, expected %s", frame, test.expected)
			}
		})
	}
}

func TestRawFramerBuffer(t *testing.T) {
	framer := framer.GetFramer("RAW")
	framer.Decode([]byte("Hello, World!"))

	if !slices.Equal(framer.Buffer(), []byte{}) {
		t.Errorf("raw framer buffer got %s, expected empty", framer.Buffer())
	}
	framer.Clear()
	if !slices.Equal(framer.Buffer(), []byte{}) {
		t.Errorf("raw framer buffer got %s, expected empty after clear", framer.Buffer())
	}
}
