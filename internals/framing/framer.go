package framing

type Framer interface {
	Frame([]byte) [][]byte
	Clear()
}
