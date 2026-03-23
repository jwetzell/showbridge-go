package processor_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
)

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

type TestProcessor struct {
}

func (p *TestProcessor) Type() string {
	return "test"
}
func (p *TestProcessor) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	return wrappedPayload, nil
}

func NewTestKVModule(id string) *TestKVModule {
	return &TestKVModule{
		id: id,
	}
}

type TestKVModule struct {
	id     string
	kvData map[string]any
}

func (m *TestKVModule) Start(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (m *TestKVModule) Stop() {}

func (m *TestKVModule) Type() string {
	return "module.test.kv"
}

func (m *TestKVModule) Id() string {
	return m.id
}

func (m *TestKVModule) Get(key string) (any, error) {
	return key, nil
}

func (m *TestKVModule) Set(key string, value any) error {
	if m.kvData == nil {
		m.kvData = make(map[string]any)
	}
	m.kvData[key] = value
	return nil
}
func NewTestDBModule(id string) *TestDBModule {
	return &TestDBModule{
		id: id,
	}
}

type TestDBModule struct {
	id string
	db *sql.DB
}

func (m *TestDBModule) Start(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (m *TestDBModule) Database() *sql.DB {
	if m.db == nil {
		db, _ := sql.Open("sqlite", ":memory:")

		db.Exec(`
		CREATE TABLE test (
			id INTEGER PRIMARY KEY,
			value TEXT
		);
		INSERT INTO test (id, value) VALUES (1, 'test-1'), (2, 'test-2');
		
		`)
		m.db = db
	}
	return m.db
}

func (m *TestDBModule) Stop() {}

func (m *TestDBModule) Type() string {
	return "module.test.db"
}

func (m *TestDBModule) Id() string {
	return m.id
}

func GetContextWithModules(ctx context.Context, modules map[string]common.Module) context.Context {
	ctx = context.WithValue(ctx, common.ModulesContextKey, modules)
	return ctx
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
