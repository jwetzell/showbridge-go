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
		return newByteSeparatorFramer([]byte{'\r'})
	case "LF":
		return newByteSeparatorFramer([]byte{'\n'})
	case "CRLF":
		return newByteSeparatorFramer([]byte{'\r', '\n'})
	case "SLIP":
		return newSlipFramer()
	case "RAW":
		return newRawFramer()
	default:
		return nil
	}
}
