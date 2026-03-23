package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
	"github.com/jwetzell/showbridge-go/internal/test"
)

func TestModuleBadRegistrationNoType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("module registration should have panicked but did not")
		}
	}()

	module.RegisterModule(module.ModuleRegistration{
		Type: "",
		New: func(config config.ModuleConfig) (common.Module, error) {
			return &test.TestModule{}, nil
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
		New: func(config config.ModuleConfig) (common.Module, error) {
			return &test.TestModule{}, nil
		},
	})

	module.RegisterModule(module.ModuleRegistration{
		Type: "module.test",
		New: func(config config.ModuleConfig) (common.Module, error) {
			return &test.TestModule{}, nil
		},
	})
}
