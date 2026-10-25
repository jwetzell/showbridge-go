package processor_test

import (
	"context"
	"testing"

	"github.com/jwetzell/showbridge-go/common"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/payload"
	"github.com/jwetzell/showbridge-go/processor"
)

type testModuleManager map[string]common.Module

func (m testModuleManager) GetModule(name string) (common.Module, bool) {
	mod, ok := m[name]
	return mod, ok
}

type TestProcessor struct {
	Config config.ProcessorConfig
}

func (p *TestProcessor) Process(ctx context.Context, wrappedPayload payload.WrappedPayload, moduleManager common.ModuleManager, inputHandler common.InputHandler) (payload.WrappedPayload, error) {
	return wrappedPayload, nil
}

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

func TestProcessorBadRegistrationNoType(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("processor registration should have panicked but did not")
		}
	}()

	processor.RegisterProcessor(processor.ProcessorRegistration{
		Type: "",
		New: func(config config.ProcessorConfig) (processor.Processor, error) {
			return &TestProcessor{}, nil
		},
	})
}

func TestProcessorBadRegistrationNoNew(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("processor registration should have panicked but did not")
		}
	}()

	processor.RegisterProcessor(processor.ProcessorRegistration{
		Type: "test",
		New:  nil,
	})
}

func TestProcessorBadRegistrationExistingType(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("processor registration should have panicked but did not")
		}
	}()

	processor.RegisterProcessor(processor.ProcessorRegistration{
		Type: "string.create",
		New: func(config config.ProcessorConfig) (processor.Processor, error) {
			return &TestProcessor{}, nil
		},
	})
}
