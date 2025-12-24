package framer_test

import (
	"slices"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/framer"
)

func TestGoodRawFramer(t *testing.T) {
	tests := []struct {
		name     string
		framer   framer.Framer
		input    []byte
		expected [][]byte
	}{
		{
			name:   "basic raw framer",
			framer: framer.NewRawFramer(),
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
