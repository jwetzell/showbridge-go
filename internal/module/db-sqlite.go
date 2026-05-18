package module

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"

	_ "modernc.org/sqlite"
)

type DbSqlite struct {
	config config.ModuleConfig
	Dsn    string
	ctx    context.Context
	router common.RouteIO
	db     *sql.DB
	logger *slog.Logger
	dbMu   sync.Mutex
	cancel context.CancelFunc
}

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "db.sqlite",
		Title: "SQLite Database",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"dsn": {
					Type:      "string",
					MinLength: new(1),
				},
			},
			Required:             []string{"dsn"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ModuleConfig) (common.Module, error) {
			params := config.Params

			dsnString, err := params.GetString("dsn")
			if err != nil {
				return nil, fmt.Errorf("db.sqlite dsn error: %w", err)
			}

			return &DbSqlite{Dsn: dsnString, config: config, logger: CreateLogger(config)}, nil
		},
	})
}

func (t *DbSqlite) Id() string {
	return t.config.Id
}

func (t *DbSqlite) Type() string {
	return t.config.Type
}

func (t *DbSqlite) Start(ctx context.Context, router common.RouteIO) error {
	t.logger.Debug("running")
	t.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	t.ctx = moduleContext
	t.cancel = cancel

	db, err := sql.Open("sqlite", t.Dsn)
	if err != nil {
		return fmt.Errorf("db.sqlite error opening database: %w", err)
	}
	t.dbMu.Lock()
	t.db = db
	t.dbMu.Unlock()
	<-t.ctx.Done()
	return nil
}

func (t *DbSqlite) Stop() {
	if t.cancel != nil {
		t.cancel()
	}
	t.dbMu.Lock()
	defer t.dbMu.Unlock()
	if t.db != nil {
		t.db.Close()
		t.db = nil
	}
	t.logger.Debug("done")
}

// TODO(jwetzell): get a database module layout that doesn't require handing the DB over
func (t *DbSqlite) Database() (*sql.DB, error) {
	t.dbMu.Lock()
	defer t.dbMu.Unlock()
	if t.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return t.db, nil
}
