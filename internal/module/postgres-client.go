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

	_ "github.com/jackc/pgx/v5/stdlib"
)

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "postgres.client",
		Title: "PostgreSQL Client",
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
				return nil, fmt.Errorf("postgres.client url error: %w", err)
			}

			return &PostgresClient{Url: urlString, config: config, logger: CreateLogger(config)}, nil
		},
	})
}

type PostgresClient struct {
	config       config.ModuleConfig
	Url          string
	ctx          context.Context
	inputHandler common.InputHandler
	db           *sql.DB
	logger       *slog.Logger
	dbMu         sync.Mutex
	cancel       context.CancelFunc
}

func (pc *PostgresClient) Start(ctx context.Context, inputHandler common.InputHandler) error {
	pc.logger.Debug("running")
	pc.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	pc.ctx = moduleContext
	pc.cancel = cancel

	db, err := sql.Open("pgx", pc.Url)
	if err != nil {
		return fmt.Errorf("postgres.client error connecting to database: %w", err)
	}
	pc.dbMu.Lock()
	pc.db = db
	pc.dbMu.Unlock()
	<-pc.ctx.Done()
	pc.logger.Debug("done")
	return nil
}

func (pc *PostgresClient) Stop() {
	if pc.cancel != nil {
		defer pc.cancel()
	}
	pc.dbMu.Lock()
	defer pc.dbMu.Unlock()
	if pc.db != nil {
		pc.db.Close()
	}
}

func (pc *PostgresClient) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	pc.dbMu.Lock()
	defer pc.dbMu.Unlock()
	if pc.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return pc.db.QueryContext(ctx, query, args...)
}
