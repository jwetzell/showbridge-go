package module_test

import (
	"context"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

type TestModule struct {
}

func (m *TestModule) Output(ctx context.Context, payload any) error {
	return nil
}

func (m *TestModule) Start(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (m *TestModule) Stop() {}

func (m *TestModule) Type() string {
	return "module.test"
}

func (m *TestModule) Id() string {
	return "test"
}

func TestModuleBadRegistrationNoType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("module registration should have panicked but did not")
		}
	}()

	module.RegisterModule(module.ModuleRegistration{
		Type: "",
		New: func(config config.ModuleConfig) (module.Module, error) {
			return &TestModule{}, nil
		},
	})
}

func TestModuleBadRegistrationNoNew(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
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
		if r := recover(); r == nil {
			t.Fatalf("processor registration should have panicked but did not")
		}
	}()

	module.RegisterModule(module.ModuleRegistration{
		Type: "module.test",
		New: func(config config.ModuleConfig) (module.Module, error) {
			return &TestModule{}, nil
		},
	})

	module.RegisterModule(module.ModuleRegistration{
		Type: "module.test",
		New: func(config config.ModuleConfig) (module.Module, error) {
			return &TestModule{}, nil
		},
	})
}
