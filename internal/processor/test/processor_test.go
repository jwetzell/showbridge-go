package processor_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"github.com/jwetzell/showbridge-go/internal/test"
)

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
			return &test.TestProcessor{}, nil
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
			return &test.TestProcessor{}, nil
		},
	})
}

func TestTestProcessor(t *testing.T) {
	testProcessor := &test.TestProcessor{
		Config: config.ProcessorConfig{
			Id:   "test-id",
			Type: "test",
		},
	}

	if testProcessor.Id() != "test-id" {
		t.Fatalf("test processor has wrong id: %s", testProcessor.Id())
	}

	if testProcessor.Type() != "test" {
		t.Fatalf("test processor has wrong type: %s", testProcessor.Type())
	}
}
