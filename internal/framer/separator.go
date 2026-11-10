package framer

import (
	"bytes"
)

type byteSeparatorFramer struct {
	buffer    []byte
	separator []byte
}

func newByteSeparatorFramer(separator []byte) *byteSeparatorFramer {
	return &byteSeparatorFramer{separator: separator, buffer: []byte{}}
}

func (bsf *byteSeparatorFramer) Decode(data []byte) [][]byte {
	messages := [][]byte{}

	bsf.buffer = append(bsf.buffer, data...)

	parts := bytes.Split(bsf.buffer, bsf.separator)

	if len(parts) > 0 {
		bsf.buffer = parts[len(parts)-1]
		messages = parts[:len(parts)-1]
	}

	return messages
}

func (bsf *byteSeparatorFramer) Encode(data []byte) []byte {
	return append(data, bsf.separator...)
}

func (bsf *byteSeparatorFramer) Clear() {
	bsf.buffer = []byte{}
}

func (bsf *byteSeparatorFramer) Buffer() []byte {
	return bsf.buffer
}
