package module

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/common"

	_ "modernc.org/sqlite"
)

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "sqlite.client",
		Title: "SQLite Client",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"dsn": {
					Title:       "Data Source Name",
					Description: "the data source name (DSN) for the SQLite database",
					Type:        "string",
					MinLength:   new(1),
				},
			},
			Required:             []string{"dsn"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ModuleConfig) (common.Module, error) {
			params := config.Params

			dsnString, err := params.GetString("dsn")
			if err != nil {
				return nil, fmt.Errorf("sqlite.client dsn error: %w", err)
			}

			return &SqliteClient{Dsn: dsnString, config: config, logger: CreateLogger(config)}, nil
		},
	})
}

type SqliteClient struct {
	config       config.ModuleConfig
	Dsn          string
	ctx          context.Context
	inputHandler common.InputHandler
	db           *sql.DB
	logger       *slog.Logger
	dbMu         sync.Mutex
	cancel       context.CancelFunc
}

func (dbs *SqliteClient) Id() string {
	return dbs.config.Id
}

func (dbs *SqliteClient) Type() string {
	return dbs.config.Type
}

func (dbs *SqliteClient) Start(ctx context.Context, inputHandler common.InputHandler) error {
	dbs.logger.Debug("running")
	dbs.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	dbs.ctx = moduleContext
	dbs.cancel = cancel

	db, err := sql.Open("sqlite", dbs.Dsn)
	if err != nil {
		return fmt.Errorf("sqlite.client error opening database: %w", err)
	}
	dbs.dbMu.Lock()
	dbs.db = db
	dbs.dbMu.Unlock()
	<-dbs.ctx.Done()
	dbs.logger.Debug("done")
	return nil
}

func (dbs *SqliteClient) Stop() {
	if dbs.cancel != nil {
		defer dbs.cancel()
	}
	dbs.dbMu.Lock()
	defer dbs.dbMu.Unlock()
	if dbs.db != nil {
		dbs.db.Close()
	}
}

func (dbs *SqliteClient) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	dbs.dbMu.Lock()
	defer dbs.dbMu.Unlock()
	if dbs.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return dbs.db.QueryContext(ctx, query, args...)
}
