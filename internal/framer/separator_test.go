package framer_test

import (
	"slices"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/framer"
)

func TestGoodSeparatorFramer(t *testing.T) {
	tests := []struct {
		name     string
		framer   framer.Framer
		input    []byte
		expected [][]byte
		buffer   []byte
	}{
		{
			name:   "new line separator",
			framer: framer.GetFramer("LF"),
			input:  []byte("Hello\nWorld\nThis is a test\n"),
			expected: [][]byte{
				[]byte("Hello"),
				[]byte("World"),
				[]byte("This is a test"),
			},
			buffer: []byte{},
		},
		{
			name:   "CR separator",
			framer: framer.GetFramer("CR"),
			input:  []byte("Hello\rWorld\rThis is a test\r"),
			expected: [][]byte{
				[]byte("Hello"),
				[]byte("World"),
				[]byte("This is a test"),
			},
			buffer: []byte{},
		},
		{
			name:   "CRLF separator",
			framer: framer.GetFramer("CRLF"),
			input:  []byte("Hello\r\nWorld\r\nThis is a test\r\n"),
			expected: [][]byte{
				[]byte("Hello"),
				[]byte("World"),
				[]byte("This is a test"),
			},
			buffer: []byte{},
		},
		{
			name:   "extra data after separator",
			framer: framer.GetFramer("CRLF"),
			input:  []byte("Hello\r\nWorld\r\nThis is a test\r\nextra"),
			expected: [][]byte{
				[]byte("Hello"),
				[]byte("World"),
				[]byte("This is a test"),
			},
			buffer: []byte("extra"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			frames := test.framer.Decode(test.input)
			if len(frames) != len(test.expected) {
				t.Errorf("separator framer got %d frames, expected %d", len(frames), len(test.expected))
			}
			for i, frame := range frames {
				if !slices.Equal(frame, test.expected[i]) {
					t.Errorf("separator framer frame %d got %s, expected %s", i, frame, test.expected[i])
				}
			}
			if !slices.Equal(test.framer.Buffer(), test.buffer) {
				t.Errorf("separator framer buffer got %s, expected %s", test.framer.Buffer(), test.buffer)
			}
		})
	}
}
