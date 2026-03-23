package test

type TestStruct struct {
	String   string
	Int      int
	Float    float64
	Bool     bool
	Data     any
	IntSlice []int
}

func (t TestStruct) GetString() string {
	return t.String
}

func (t TestStruct) GetInt() int {
	return t.Int
}

func (t TestStruct) GetFloat() float64 {
	return t.Float
}

func (t TestStruct) GetBool() bool {
	return t.Bool
}

func (t TestStruct) GetData() any {
	return t.Data
}

func (t TestStruct) GetIntSlice() []int {
	return t.IntSlice
}

func (t TestStruct) Void() {}

func (t TestStruct) MultipleReturnValues() (string, int) {
	return t.String, t.Int
}
