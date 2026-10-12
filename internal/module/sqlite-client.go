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

func (sc *SqliteClient) Id() string {
	return sc.config.Id
}

func (sc *SqliteClient) Type() string {
	return sc.config.Type
}

func (sc *SqliteClient) Start(ctx context.Context, inputHandler common.InputHandler) error {
	sc.logger.Debug("running")
	sc.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	sc.ctx = moduleContext
	sc.cancel = cancel

	db, err := sql.Open("sqlite", sc.Dsn)
	if err != nil {
		return fmt.Errorf("sqlite.client error opening database: %w", err)
	}
	sc.dbMu.Lock()
	sc.db = db
	sc.dbMu.Unlock()
	<-sc.ctx.Done()
	sc.logger.Debug("done")
	return nil
}

func (sc *SqliteClient) Stop() {
	if sc.cancel != nil {
		defer sc.cancel()
	}
	sc.dbMu.Lock()
	defer sc.dbMu.Unlock()
	if sc.db != nil {
		sc.db.Close()
	}
}

func (sc *SqliteClient) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	sc.dbMu.Lock()
	defer sc.dbMu.Unlock()
	if sc.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return sc.db.QueryContext(ctx, query, args...)
}
