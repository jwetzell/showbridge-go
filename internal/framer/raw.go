package framer

type rawFramer struct{}

func newRawFramer() *rawFramer {
	return &rawFramer{}
}

func (rf *rawFramer) Decode(data []byte) [][]byte {
	return [][]byte{data}
}

func (rf *rawFramer) Encode(data []byte) []byte {
	return data
}

func (rf *rawFramer) Clear() {
	// NOTE(jwetzell): no internal state to clear
}

func (rf *rawFramer) Buffer() []byte {
	return []byte{}
}
