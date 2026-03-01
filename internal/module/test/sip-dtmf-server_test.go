package module_test

import (
	"testing"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
)

func TestSIPDTMFServerFromRegistry(t *testing.T) {
	registration, ok := module.ModuleRegistry["sip.dtmf.server"]
	if !ok {
		t.Fatalf("sip.dtmf.server module not registered")
	}

	moduleInstance, err := registration.New(config.ModuleConfig{
		Id:   "test",
		Type: "sip.dtmf.server",
		Params: map[string]any{
			"separator": "#",
		},
	})

	if err != nil {
		t.Fatalf("failed to create sip.dtmf.server module: %s", err)
	}

	if moduleInstance.Id() != "test" {
		t.Fatalf("sip.dtmf.server module has wrong id: %s", moduleInstance.Id())
	}

	if moduleInstance.Type() != "sip.dtmf.server" {
		t.Fatalf("sip.dtmf.server module has wrong type: %s", moduleInstance.Type())
	}
}

func TestBadSIPDTMFServer(t *testing.T) {
	tests := []struct {
		name        string
		params      map[string]any
		errorString string
	}{
		{
			name:        "no separator param",
			params:      map[string]any{},
			errorString: "sip.dtmf.server separator error: not found",
		},
		{
			name: "non-string separator param",
			params: map[string]any{
				"separator": 123,
			},
			errorString: "sip.dtmf.server separator error: not a string",
		},
		{
			name: "non-number port param",
			params: map[string]any{
				"separator": "#",
				"port":      "8000",
			},
			errorString: "sip.dtmf.server port error: not a number",
		},
		{
			name: "non-string ip param",
			params: map[string]any{
				"separator": "#",
				"ip":        123,
			},
			errorString: "sip.dtmf.server ip error: not a string",
		},
		{
			name: "non-string transport param",
			params: map[string]any{
				"separator": "#",
				"transport": 123,
			},
			errorString: "sip.dtmf.server transport error: not a string",
		},
		{
			name: "non-string userAgent param",
			params: map[string]any{
				"separator": "#",
				"userAgent": 123,
			},
			errorString: "sip.dtmf.server userAgent error: not a string",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := module.ModuleRegistry["sip.dtmf.server"]
			if !ok {
				t.Fatalf("sip.dtmf.server module not registered")
			}

			moduleInstance, err := registration.New(config.ModuleConfig{
				Id:     "test",
				Type:   "sip.dtmf.server",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("sip.dtmf.server got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			err = moduleInstance.Start(t.Context())

			if err == nil {
				t.Fatalf("sip.dtmf.server expected to fail")
			}

			if err.Error() != test.errorString {
				t.Fatalf("sip.dtmf.server got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
