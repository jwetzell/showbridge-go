package module

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "db.mysql",
		Title: "MySQL Database",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"dsn": {
					Title:       "Database DSN",
					Description: "the connection DSN for the MySQL database",
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
				return nil, fmt.Errorf("db.mysql dsn error: %w", err)
			}

			return &DbMysql{Dsn: dsnString, config: config, logger: CreateLogger(config)}, nil
		},
	})
}

type DbMysql struct {
	config       config.ModuleConfig
	Dsn          string
	ctx          context.Context
	inputHandler common.InputHandler
	db           *sql.DB
	logger       *slog.Logger
	dbMu         sync.Mutex
	cancel       context.CancelFunc
}

func (dbs *DbMysql) Id() string {
	return dbs.config.Id
}

func (dbs *DbMysql) Type() string {
	return dbs.config.Type
}

func (dbs *DbMysql) Start(ctx context.Context, inputHandler common.InputHandler) error {
	dbs.logger.Debug("running")
	dbs.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	dbs.ctx = moduleContext
	dbs.cancel = cancel

	db, err := sql.Open("mysql", dbs.Dsn)
	if err != nil {
		return fmt.Errorf("db.mysql error connecting to database: %w", err)
	}

	// TODO(jwetzell): make configurable
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	dbs.dbMu.Lock()
	dbs.db = db
	dbs.dbMu.Unlock()
	<-dbs.ctx.Done()
	dbs.logger.Debug("done")
	return nil
}

func (dbs *DbMysql) Stop() {
	if dbs.cancel != nil {
		defer dbs.cancel()
	}
	dbs.dbMu.Lock()
	defer dbs.dbMu.Unlock()
	if dbs.db != nil {
		dbs.db.Close()
	}
}

func (dbs *DbMysql) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	dbs.dbMu.Lock()
	defer dbs.dbMu.Unlock()
	if dbs.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return dbs.db.QueryContext(ctx, query, args...)
}
