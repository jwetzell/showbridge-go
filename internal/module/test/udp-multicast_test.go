package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestUDPMulticastFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["net.udp.multicast"]
	if !ok {
		t.Fatalf("udp.multicast module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "net.udp.multicast",
		Params: map[string]any{
			"ip":   "236.10.10.10",
			"port": 56565.0,
		},
	})

	if err != nil {
		t.Fatalf("failed to create net.udp.multicast module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("net.udp.multicast module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "net.udp.multicast" {
		t.Fatalf("net.udp.multicast module has wrong type: %s", moduleInstance.Type())
	}
}

func TestBadUDPMulticast(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name: "no port param",
			params: map[string]any{
				"ip": "localhost",
			},
			errorString: "net.udp.multicast requires a port parameter",
		},
		{
			name: "non-number port param",
			params: map[string]any{
				"ip":   "localhost",
				"port": "8000",
			},
			errorString: "net.udp.multicast port must be a number",
		},
		{
			name: "no ip param",
			params: map[string]any{
				"port": 8000.0,
			},
			errorString: "net.udp.multicast requires an ip parameter",
		},
		{
			name: "non-string ip param",
			params: map[string]any{
				"ip":   123,
				"port": 8000.0,
			},
			errorString: "net.udp.multicast ip must be a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.ModuleRegistry["net.udp.multicast"]
			if !ok {
				t.Fatalf("net.udp.multicast module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "net.udp.multicast",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("net.udp.multicast got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context())

			if err == nil {
				t.Fatalf("net.udp.multicast expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("net.udp.multicast got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
