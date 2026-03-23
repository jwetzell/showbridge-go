package processor_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
	_ "modernc.org/sqlite"
)

func TestDbQueryFromRegistry(t *testing.T) {
	registration, ok := processor.ProcessorRegistry["db.query"]
	if !ok {
		t.Fatalf("db.query processor not registered")
	}

	processorInstance, err := registration.New(config.ProcessorConfig{
		Type: "db.query",
		Params: map[string]any{
			"module": "test",
			"query":  "SELECT sqlite_version();",
		},
	})
	if err != nil {
		t.Fatalf("failed to create db.query processor: %s", err)
	}

	if processorInstance.Type() != "db.query" {
		t.Fatalf("db.query processor has wrong type: %s", processorInstance.Type())
	}

	payload := "hello"
	expected := map[string]any{"sqlite_version()": "3.51.3"}

	got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(GetContextWithModules(
		t.Context(),
		map[string]common.Module{
			"test": NewTestDBModule("test"),
		},
	), payload))
	if err != nil {
		t.Fatalf("db.query processing failed: %s", err)
	}

	if !reflect.DeepEqual(got.Payload, expected) {
		t.Fatalf("db.query got %+v, expected %+v", got.Payload, expected)
	}
}

func TestGoodDbQuery(t *testing.T) {

	tests := []struct {
		name     string
		params   map[string]any
		payload  any
		expected any
	}{
		{
			name: "basic query",
			params: map[string]any{
				"module": "test",
				"query":  "select value from test where id = 1;",
			},
			payload:  "",
			expected: map[string]any{"value": "test-1"},
		},
		{
			name: "template query",
			params: map[string]any{
				"module": "test",
				"query":  "select value from test where id = {{.Payload}};",
			},
			payload:  "1",
			expected: map[string]any{"value": "test-1"},
		},
		{
			name: "multiple rows",
			params: map[string]any{
				"module": "test",
				"query":  "select * from test;",
			},
			payload: "",
			expected: []map[string]any{
				{"id": int64(1), "value": "test-1"},
				{"id": int64(2), "value": "test-2"},
			},
		},
		{
			name: "no rows",
			params: map[string]any{
				"module": "test",
				"query":  "select * from test where id = -1;",
			},
			payload:  "",
			expected: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registration, ok := processor.ProcessorRegistry["db.query"]
			if !ok {
				t.Fatalf("db.query processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "db.query",
				Params: test.params,
			})

			if err != nil {
				t.Fatalf("db.query failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(GetContextWithModules(
				t.Context(),
				map[string]common.Module{
					"test": NewTestDBModule("test"),
				},
			), test.payload))

			if err != nil {
				t.Fatalf("db.query processing failed: %s", err)
			}

			if !reflect.DeepEqual(got.Payload, test.expected) {
				t.Fatalf("db.query got payload: %+v, expected %+v", got.Payload, test.expected)
			}
		})
	}
}

func TestBadDbQuery(t *testing.T) {
	tests := []struct {
		name              string
		params            map[string]any
		payload           any
		wrappedPayloadCtx context.Context
		errorString       string
	}{
		{
			name:    "no module param",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"query": "SELECT sqlite_version();",
			},
			wrappedPayloadCtx: GetContextWithModules(t.Context(), map[string]common.Module{
				"test": NewTestDBModule("test"),
			}),
			errorString: "db.query module error: not found",
		},
		{
			name:    "non string module",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": 1,
				"query":  "SELECT sqlite_version();",
			},
			wrappedPayloadCtx: GetContextWithModules(t.Context(), map[string]common.Module{
				"test": NewTestDBModule("test"),
			}),
			errorString: "db.query module error: not a string",
		},
		{
			name:    "no query param",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
			},
			wrappedPayloadCtx: GetContextWithModules(t.Context(), map[string]common.Module{
				"test": NewTestDBModule("test"),
			}),
			errorString: "db.query query error: not found",
		},
		{
			name:    "non string query",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"query":  1,
			},
			wrappedPayloadCtx: GetContextWithModules(t.Context(), map[string]common.Module{
				"test": NewTestDBModule("test"),
			}),
			errorString: "db.query query error: not a string",
		},
		{
			name:    "query template syntax error",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"query":  "select * from {{",
			},
			wrappedPayloadCtx: GetContextWithModules(t.Context(), map[string]common.Module{
				"test": NewTestDBModule("test"),
			}),
			errorString: "template: query:1: unclosed action",
		},
		{
			name:    "query template error",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"query":  "select * from {{.Data}}",
			},
			wrappedPayloadCtx: GetContextWithModules(t.Context(), map[string]common.Module{
				"test": NewTestDBModule("test"),
			}),
			errorString: "template: query:1:16: executing \"query\" at <.Data>: can't evaluate field Data in type common.WrappedPayload",
		},
		{
			name:    "query error",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"query":  "select * from asdf;",
			},
			wrappedPayloadCtx: GetContextWithModules(t.Context(), map[string]common.Module{
				"test": NewTestDBModule("test"),
			}),
			errorString: "db.query error executing query: SQL logic error: no such table: asdf (1)",
		},
		{
			name:    "no modules in context",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"query":  "select * from test;",
			},
			wrappedPayloadCtx: t.Context(),
			errorString:       "db.query wrapped payload has no modules",
		},
		{
			name:    "module not found in context",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"query":  "select * from test;",
			},
			wrappedPayloadCtx: GetContextWithModules(t.Context(), map[string]common.Module{}),
			errorString:       "db.query unable to find module with id: test",
		},
		{
			name:    "module not found in context",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"query":  "select * from test;",
			},
			wrappedPayloadCtx: GetContextWithModules(t.Context(), map[string]common.Module{}),
			errorString:       "db.query unable to find module with id: test",
		},
		{
			name:    "module not a DatabseModule",
			payload: TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"query":  "select * from test;",
			},
			wrappedPayloadCtx: GetContextWithModules(t.Context(), map[string]common.Module{
				"test": NewTestKVModule("test"),
			}),
			errorString: "db.query module with id test is not a DatabaseModule",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.ProcessorRegistry["db.query"]
			if !ok {
				t.Fatalf("db.query processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "db.query",
				Params: test.params,
			})

			if err != nil {
				if test.errorString != err.Error() {
					t.Fatalf("db.query got error '%s', expected '%s'", err.Error(), test.errorString)
				}
				return
			}

			got, err := processorInstance.Process(t.Context(), common.GetWrappedPayload(test.wrappedPayloadCtx, test.payload))

			if err == nil {
				t.Fatalf("db.query expected to fail but got payload: %+v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("db.query got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
