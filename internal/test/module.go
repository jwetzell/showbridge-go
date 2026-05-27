package test

import (
	"context"
	"database/sql"

	"github.com/jwetzell/showbridge-go/internal/common"
	_ "modernc.org/sqlite"
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

func (m *TestModule) Type() string {
	return "test.plain"
}

func (m *TestModule) Id() string {
	return "test"
}

func NewTestOutputModule(id string) *TestOutputModule {
	return &TestOutputModule{
		id: id,
	}
}

type TestOutputModule struct {
	id string
}

func (m *TestOutputModule) Start(ctx context.Context, inputHandler common.InputHandler) error {
	<-ctx.Done()
	return nil
}

func (m *TestOutputModule) Output(ctx context.Context, payload any) error {
	return nil
}

func (m *TestOutputModule) Stop() {}

func (m *TestOutputModule) Type() string {
	return "test.output"
}

func (m *TestOutputModule) Id() string {
	return m.id
}

func NewTestKVModule(id string, presetValues map[string]any) *TestKVModule {
	return &TestKVModule{
		id:     id,
		kvData: presetValues,
	}
}

type TestKVModule struct {
	id     string
	kvData map[string]any
}

func (m *TestKVModule) Start(ctx context.Context, inputHandler common.InputHandler) error {
	<-ctx.Done()
	return nil
}

func (m *TestKVModule) Stop() {}

func (m *TestKVModule) Type() string {
	return "test.kv"
}

func (m *TestKVModule) Id() string {
	return m.id
}

func (m *TestKVModule) Get(ctx context.Context, key string) (any, error) {
	if m.kvData == nil {
		return nil, nil
	}
	value, ok := m.kvData[key]
	if !ok {
		return nil, nil
	}
	return value, nil
}

func (m *TestKVModule) Set(ctx context.Context, key string, value any) error {
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

func (m *TestDBModule) Start(ctx context.Context, inputHandler common.InputHandler) error {
	<-ctx.Done()
	return nil
}

func (m *TestDBModule) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if m.db == nil {
		db, err := sql.Open("sqlite", ":memory:")
		if err != nil {
			return nil, err
		}

		_, err = db.Exec(`
		CREATE TABLE test (
			id INTEGER PRIMARY KEY,
			value TEXT
		);
		INSERT INTO test (id, value) VALUES (1, 'test-1'), (2, 'test-2');
		
		`)
		if err != nil {
			return nil, err
		}
		m.db = db
	}
	return m.db.QueryContext(ctx, query, args...)
}

func (m *TestDBModule) Stop() {}

func (m *TestDBModule) Type() string {
	return "test.db"
}

func (m *TestDBModule) Id() string {
	return m.id
}

func NewTestPubSubModule(id string) *TestPubSubModule {
	return &TestPubSubModule{
		id: id,
	}
}

type TestPubSubModule struct {
	id string
}

func (m *TestPubSubModule) Start(ctx context.Context, inputHandler common.InputHandler) error {
	<-ctx.Done()
	return nil
}

func (m *TestPubSubModule) Publish(ctx context.Context, topic string, payload any) error {
	return nil
}

func (m *TestPubSubModule) Stop() {}

func (m *TestPubSubModule) Type() string {
	return "test.pubsub"
}

func (m *TestPubSubModule) Id() string {
	return m.id
}
