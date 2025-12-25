package framer

type Framer interface {
	Decode([]byte) [][]byte
	Encode([]byte) []byte
	Clear()
	Buffer() []byte
}

func GetFramer(framingType string) Framer {
	switch framingType {
	case "CR":
		return NewByteSeparatorFramer([]byte{'\r'})
	case "LF":
		return NewByteSeparatorFramer([]byte{'\n'})
	case "CRLF":
		return NewByteSeparatorFramer([]byte{'\r', '\n'})
	case "SLIP":
		return NewSlipFramer()
	case "RAW":
		return NewRawFramer()
	default:
		return nil
	}
}
