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

	_ "github.com/jackc/pgx/v5/stdlib"
)

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "db.postgres",
		Title: "PostgreSQL Database",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"url": {
					Title:       "Database URL",
					Description: "the connection URL for the PostgreSQL database",
					Type:        "string",
					MinLength:   new(1),
				},
			},
			Required:             []string{"url"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ModuleConfig) (common.Module, error) {
			params := config.Params

			urlString, err := params.GetString("url")
			if err != nil {
				return nil, fmt.Errorf("db.postgres url error: %w", err)
			}

			return &DbPostgres{Url: urlString, config: config, logger: CreateLogger(config)}, nil
		},
	})
}

type DbPostgres struct {
	config       config.ModuleConfig
	Url          string
	ctx          context.Context
	inputHandler common.InputHandler
	db           *sql.DB
	logger       *slog.Logger
	dbMu         sync.Mutex
	cancel       context.CancelFunc
}

func (dbs *DbPostgres) Id() string {
	return dbs.config.Id
}

func (dbs *DbPostgres) Type() string {
	return dbs.config.Type
}

func (dbs *DbPostgres) Start(ctx context.Context, inputHandler common.InputHandler) error {
	dbs.logger.Debug("running")
	dbs.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	dbs.ctx = moduleContext
	dbs.cancel = cancel

	db, err := sql.Open("pgx", dbs.Url)
	if err != nil {
		return fmt.Errorf("db.postgres error connecting to database: %w", err)
	}
	dbs.dbMu.Lock()
	dbs.db = db
	dbs.dbMu.Unlock()
	<-dbs.ctx.Done()
	dbs.logger.Debug("done")
	return nil
}

func (dbs *DbPostgres) Stop() {
	if dbs.cancel != nil {
		defer dbs.cancel()
	}
	dbs.dbMu.Lock()
	defer dbs.dbMu.Unlock()
	if dbs.db != nil {
		dbs.db.Close()
	}
}

func (dbs *DbPostgres) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	dbs.dbMu.Lock()
	defer dbs.dbMu.Unlock()
	if dbs.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return dbs.db.QueryContext(ctx, query, args...)
}
