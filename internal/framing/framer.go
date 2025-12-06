package framing

import (
	"fmt"
)

type Framer interface {
	Decode([]byte) [][]byte
	Encode([]byte) []byte
	Clear()
}

func GetFramer(framingType string) (Framer, error) {
	switch framingType {
	case "CR":
		return NewByteSeparatorFramer([]byte{'\r'}), nil
	case "LF":
		return NewByteSeparatorFramer([]byte{'\n'}), nil
	case "CRLF":
		return NewByteSeparatorFramer([]byte{'\r', '\n'}), nil
	case "SLIP":
		return NewSlipFramer(), nil
	default:
		return nil, fmt.Errorf("unknown framing method: %s", framingType)
	}
}
