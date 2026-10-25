package framer

import "io"

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

type FramedWriter struct {
	writer io.Writer
	framer Framer
}

func NewFramedWrite(writer io.Writer, framer Framer) *FramedWriter {
	return &FramedWriter{
		writer: writer,
		framer: framer,
	}
}

func (fw *FramedWriter) Write(data []byte) (int, error) {
	encoded := fw.framer.Encode(data)
	return fw.writer.Write(encoded)
}
