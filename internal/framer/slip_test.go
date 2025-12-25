package framer_test

import (
	"slices"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/framer"
)

func TestGoodSLIPFramerDecode(t *testing.T) {
	tests := []struct {
		name     string
		framer   framer.Framer
		input    []byte
		expected [][]byte
		buffer   []byte
	}{
		{
			name:   "OSC SLIP messages",
			framer: framer.GetFramer("SLIP"),
			input:  []byte{0xc0, 0x2f, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x00, 0x00, 0x2c, 0x00, 0x00, 0x00, 0xc0},
			expected: [][]byte{
				{0x2f, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x00, 0x00, 0x2c, 0x00, 0x00, 0x00},
			},
			buffer: []byte{},
		},
		{
			name:     "SLIP decode escaped end",
			framer:   framer.GetFramer("SLIP"),
			expected: [][]byte{{0xc0}},
			input:    []byte{0xc0, 0xdb, 0xdc, 0xc0},
			buffer:   []byte{},
		},
		{
			name:     "SLIP decode escaped escape",
			framer:   framer.GetFramer("SLIP"),
			expected: [][]byte{{0xdb}},
			input:    []byte{0xc0, 0xdb, 0xdd, 0xc0},
			buffer:   []byte{},
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

func TestGoodSLIPFramerEncode(t *testing.T) {
	tests := []struct {
		name     string
		framer   framer.Framer
		input    []byte
		expected []byte
	}{
		{
			name:   "OSC SLIP messages",
			framer: framer.GetFramer("SLIP"),
			input: []byte{
				0x2f, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x00, 0x00, 0x2c, 0x00, 0x00, 0x00,
			},
			expected: []byte{0xc0, 0x2f, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x00, 0x00, 0x2c, 0x00, 0x00, 0x00, 0xc0},
		},
		{
			name:     "SLIP encode end",
			framer:   framer.GetFramer("SLIP"),
			input:    []byte{0xc0},
			expected: []byte{0xc0, 0xdb, 0xdc, 0xc0},
		},
		{
			name:     "SLIP encode esc",
			framer:   framer.GetFramer("SLIP"),
			input:    []byte{0xdb},
			expected: []byte{0xc0, 0xdb, 0xdd, 0xc0},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			frame := test.framer.Encode(test.input)
			if !slices.Equal(frame, test.expected) {
				t.Errorf("SLIP framer frame got %s, expected %s", frame, test.expected)
			}
		})
	}
}

func TestSlipFramerBuffer(t *testing.T) {
	framer := framer.GetFramer("SLIP")
	framer.Decode([]byte{0xc0, 0x2f, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x00, 0x00, 0x2c, 0x00, 0x00, 0x00, 0xc0, 0xc0, 0x45})
	if !slices.Equal(framer.Buffer(), []byte{0x45}) {
		t.Errorf("SLIP framer buffer got %s, expected %s", framer.Buffer(), []byte{0x45})
	}
	framer.Clear()
	if !slices.Equal(framer.Buffer(), []byte{}) {
		t.Errorf("SLIP framer buffer got %s, expected empty slice", framer.Buffer())
	}
}
