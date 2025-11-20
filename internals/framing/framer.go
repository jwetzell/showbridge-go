package framing

type Framer interface {
	Decode([]byte) [][]byte
	Encode([]byte) []byte
	Clear()
}
