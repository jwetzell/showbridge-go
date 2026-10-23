package processor_test

import (
	"reflect"
	"testing"

	"github.com/jwetzell/showbridge-go/common"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/test"
	"github.com/jwetzell/showbridge-go/processor"
	_ "modernc.org/sqlite"
)

func TestGoodDbQuery(t *testing.T) {

	testCases := []struct {
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
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			registration, ok := processor.GetProcessorRegistration("db.query")
			if !ok {
				t.Fatalf("db.query processor not registered")
			}

			processorInstance, err := registration.New(config.ProcessorConfig{
				Type:   "db.query",
				Params: testCase.params,
			})

			if err != nil {
				t.Fatalf("db.query failed to create processor: %s", err)
			}

			got, err := processorInstance.Process(t.Context(), common.WrappedPayload{
				Payload: testCase.payload,
			}, testModuleManager{
				"test": test.NewTestDBModule("test"),
			}, nil)

			if err != nil {
				t.Fatalf("db.query processing failed: %s", err)
			}

			if !reflect.DeepEqual(got.Payload, testCase.expected) {
				t.Fatalf("db.query got payload: %+v, expected %+v", got.Payload, testCase.expected)
			}
		})
	}
}

func TestBadDbQuery(t *testing.T) {
	tests := []struct {
		name          string
		params        map[string]any
		payload       any
		moduleManager common.ModuleManager
		errorString   string
	}{
		{
			name:    "no module param",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"query": "SELECT sqlite_version();",
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": test.NewTestDBModule("test"),
			}),
			errorString: "db.query module error: not found",
		},
		{
			name:    "non string module",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": 1,
				"query":  "SELECT sqlite_version();",
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": test.NewTestDBModule("test"),
			}),
			errorString: "db.query module error: not a string",
		},
		{
			name:    "no query param",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": test.NewTestDBModule("test"),
			}),
			errorString: "db.query query error: not found",
		},
		{
			name:    "non string query",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"query":  1,
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": test.NewTestDBModule("test"),
			}),
			errorString: "db.query query error: not a string",
		},
		{
			name:    "query template syntax error",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"query":  "select * from {{",
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": test.NewTestDBModule("test"),
			}),
			errorString: "template: query:1: unclosed action",
		},
		{
			name:    "query template error",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"query":  "select * from {{.Data}}",
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": test.NewTestDBModule("test"),
			}),
			errorString: "template: query:1:16: executing \"query\" at <.Data>: can't evaluate field Data in type common.WrappedPayload",
		},
		{
			name:    "query error",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"query":  "select * from asdf;",
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": test.NewTestDBModule("test"),
			}),
			errorString: "db.query error executing query: SQL logic error: no such table: asdf (1)",
		},
		{
			name:    "no modules",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"query":  "select * from test;",
			},
			moduleManager: nil,
			errorString:   "db.query module manager is nil",
		},
		{
			name:    "module not found",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"query":  "select * from test;",
			},
			moduleManager: testModuleManager(map[string]common.Module{}),
			errorString:   "db.query unable to find module with id: test",
		},
		{
			name:    "module not a DatabaseModule",
			payload: test.TestStruct{Data: "hello"},
			params: map[string]any{
				"module": "test",
				"query":  "select * from test;",
			},
			moduleManager: testModuleManager(map[string]common.Module{
				"test": test.NewTestKVModule("test", nil),
			}),
			errorString: "db.query module with id test is not a DatabaseModule",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			registration, ok := processor.GetProcessorRegistration("db.query")
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

			got, err := processorInstance.Process(t.Context(), common.WrappedPayload{
				Payload: test.payload,
			}, test.moduleManager, nil)

			if err == nil {
				t.Fatalf("db.query expected to fail but got payload: %+v", got)
			}

			if err.Error() != test.errorString {
				t.Fatalf("db.query got error '%s', expected '%s'", err.Error(), test.errorString)
			}
		})
	}
}
