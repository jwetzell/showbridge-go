package framer

import (
	"bytes"

	"github.com/jwetzell/slip-go"
)

type slipFramer struct {
	buffer []byte
}

func newSlipFramer() *slipFramer {
	return &slipFramer{buffer: []byte{}}
}

func (sf *slipFramer) Decode(data []byte) [][]byte {

	messages := [][]byte{}

	END := byte(0xc0)
	ESC := byte(0xdb)
	ESC_END := byte(0xdc)
	ESC_ESC := byte(0xdd)

	escapeNext := false
	for _, packetByte := range data {

		if packetByte == ESC {
			escapeNext = true
			continue
		}

		if escapeNext {
			switch packetByte {
			case ESC_END:
				sf.buffer = append(sf.buffer, END)
			case ESC_ESC:
				sf.buffer = append(sf.buffer, ESC)
			}
			escapeNext = false
		} else if packetByte == END {
			if len(sf.buffer) == 0 {
				// opening END byte, can discard
				continue
			} else {
				message := sf.buffer
				messages = append(messages, message)
			}
			sf.buffer = []byte{}
		} else {
			sf.buffer = append(sf.buffer, packetByte)
		}
	}

	return messages
}

func (sf *slipFramer) Encode(data []byte) []byte {
	encodedBuffer := bytes.NewBuffer(nil)
	w := slip.NewWriter(encodedBuffer)
	_, _ = w.Write(data)
	return encodedBuffer.Bytes()
}

func (sf *slipFramer) Clear() {
	sf.buffer = []byte{}
}

func (sf *slipFramer) Buffer() []byte {
	return sf.buffer
}
