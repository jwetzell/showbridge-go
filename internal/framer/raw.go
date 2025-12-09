package framer

type RawFramer struct{}

func NewRawFramer() *RawFramer {
	return &RawFramer{}
}

func (rf *RawFramer) Decode(data []byte) [][]byte {
	return [][]byte{data}
}

func (rf *RawFramer) Encode(data []byte) []byte {
	return data
}

func (rf *RawFramer) Clear() {
	// NOTE(jwetzell): no internal state to clear
}

func (rf *RawFramer) Buffer() []byte {
	return []byte{}
}
