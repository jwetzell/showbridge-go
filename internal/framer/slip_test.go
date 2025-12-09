package framer_test

import (
	"slices"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/framer"
)

func TestGoodSLIPFramer(t *testing.T) {
	tests := []struct {
		name     string
		framer   framer.Framer
		input    []byte
		expected [][]byte
		buffer   []byte
	}{
		{
			name:   "OSC SLIP messages",
			framer: framer.NewSlipFramer(),
			input:  []byte{0xc0, 0x2f, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x00, 0x00, 0x2c, 0x00, 0x00, 0x00, 0xc0},
			expected: [][]byte{
				{0x2f, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x00, 0x00, 0x2c, 0x00, 0x00, 0x00},
			},
			buffer: []byte{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			frames := test.framer.Decode(test.input)
			if len(frames) != len(test.expected) {
				t.Errorf("SLIP framer got %d frames, expected %d", len(frames), len(test.expected))
			}
			for i, frame := range frames {
				if !slices.Equal(frame, test.expected[i]) {
					t.Errorf("SLIP framer frame %d got %s, expected %s", i, frame, test.expected[i])
				}
			}
			if !slices.Equal(test.framer.Buffer(), test.buffer) {
				t.Errorf("SLIP framer buffer got %s, expected %s", test.framer.Buffer(), test.buffer)
			}
		})
	}
}
