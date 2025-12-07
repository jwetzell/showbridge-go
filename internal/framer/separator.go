package framer

import (
	"bytes"
)

type ByteSeparatorFramer struct {
	buffer    []byte
	separator []byte
}

func NewByteSeparatorFramer(separator []byte) *ByteSeparatorFramer {
	return &ByteSeparatorFramer{separator: separator, buffer: []byte{}}
}

func (bsf *ByteSeparatorFramer) Decode(data []byte) [][]byte {
	messages := [][]byte{}

	bsf.buffer = append(bsf.buffer, data...)

	parts := bytes.Split(bsf.buffer, bsf.separator)

	if len(parts) > 0 {
		bsf.buffer = parts[len(parts)-1]
		messages = parts[:len(parts)-1]
	}

	return messages
}

func (bsf *ByteSeparatorFramer) Encode(data []byte) []byte {
	return append(data, bsf.separator...)
}

func (bsf *ByteSeparatorFramer) Clear() {
	bsf.buffer = []byte{}
}
