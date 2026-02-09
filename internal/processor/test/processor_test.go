package processor_test

import (
	"context"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

type TestProcessor struct {
}

func (p *TestProcessor) Type() string {
	return "test"
}
func (p *TestProcessor) Process(ctx context.Context, input any) (any, error) {
	return input, nil
}

func TestProcessorBadRegistrationNoType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
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
		if r := recover(); r == nil {
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
		if r := recover(); r == nil {
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
