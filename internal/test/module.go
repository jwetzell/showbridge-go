package test

import (
	"context"
	"database/sql"

	"github.com/jwetzell/showbridge-go/internal/common"
	_ "modernc.org/sqlite"
)

type TestModule struct {
}

func (m *TestModule) Start(ctx context.Context, router common.RouteIO) error {
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

func NewTestKVModule(id string) *TestKVModule {
	return &TestKVModule{
		id: id,
	}
}

type TestKVModule struct {
	id     string
	kvData map[string]any
}

func (m *TestKVModule) Start(ctx context.Context, router common.RouteIO) error {
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

func (m *TestDBModule) Start(ctx context.Context, router common.RouteIO) error {
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
