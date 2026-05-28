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

type DbSqlite struct {
	config       config.ModuleConfig
	Dsn          string
	ctx          context.Context
	inputHandler common.InputHandler
	db           *sql.DB
	logger       *slog.Logger
	dbMu         sync.Mutex
	cancel       context.CancelFunc
}

func (dbs *DbSqlite) Id() string {
	return dbs.config.Id
}

func (dbs *DbSqlite) Type() string {
	return dbs.config.Type
}

func (dbs *DbSqlite) Start(ctx context.Context, inputHandler common.InputHandler) error {
	dbs.logger.Debug("running")
	dbs.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	dbs.ctx = moduleContext
	dbs.cancel = cancel

	db, err := sql.Open("sqlite", dbs.Dsn)
	if err != nil {
		return fmt.Errorf("db.sqlite error opening database: %w", err)
	}
	dbs.dbMu.Lock()
	dbs.db = db
	dbs.dbMu.Unlock()
	<-dbs.ctx.Done()
	dbs.logger.Debug("done")
	return nil
}

func (dbs *DbSqlite) Stop() {
	if dbs.cancel != nil {
		defer dbs.cancel()
	}
	dbs.dbMu.Lock()
	defer dbs.dbMu.Unlock()
	if dbs.db != nil {
		dbs.db.Close()
	}
}

func (dbs *DbSqlite) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	dbs.dbMu.Lock()
	defer dbs.dbMu.Unlock()
	if dbs.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return dbs.db.QueryContext(ctx, query, args...)
}
