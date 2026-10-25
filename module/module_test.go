package module_test

import (
	"context"
	"testing"

	"github.com/jwetzell/showbridge-go/common"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/module"
)

func NewTestModule(id string) *TestModule {
	return &TestModule{
		id: id,
	}
}

type TestModule struct {
	id string
}

func (m *TestModule) Start(ctx context.Context, inputHandler common.InputHandler) error {
	<-ctx.Done()
	return nil
}

func (m *TestModule) Stop() {}

func TestModuleBadRegistrationNoType(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("module registration should have panicked but did not")
		}
	}()

	module.RegisterModule(module.ModuleRegistration{
		Type: "",
		New: func(config config.ModuleConfig) (common.Module, error) {
			return NewTestModule("test"), nil
		},
	})
}

func TestModuleBadRegistrationNoNew(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("processor registration should have panicked but did not")
		}
	}()

	module.RegisterModule(module.ModuleRegistration{
		Type: "module.test",
		New:  nil,
	})
}

func TestModuleBadRegistrationExistingType(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("processor registration should have panicked but did not")
		}
	}()

	module.RegisterModule(module.ModuleRegistration{
		Type: "module.test",
		New: func(config config.ModuleConfig) (common.Module, error) {
			return NewTestModule("test"), nil
		},
	})

	module.RegisterModule(module.ModuleRegistration{
		Type: "module.test",
		New: func(config config.ModuleConfig) (common.Module, error) {
			return NewTestModule("test"), nil
		},
	})
}
