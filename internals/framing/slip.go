package framing

type SlipFramer struct {
	buffer []byte
}

func NewSlipFramer() *SlipFramer {
	return &SlipFramer{buffer: []byte{}}
}

func (sf *SlipFramer) Frame(data []byte) [][]byte {
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
			if packetByte == ESC_END {
				sf.buffer = append(sf.buffer, END)
			} else if packetByte == ESC_ESC {
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

func (sf *SlipFramer) Clear() {
	sf.buffer = []byte{}
}
