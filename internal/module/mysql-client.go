package module

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/common"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "mysql.client",
		Title: "MySQL Client",
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
				return nil, fmt.Errorf("mysql.client dsn error: %w", err)
			}

			return &MySQLClient{Dsn: dsnString, config: config, logger: CreateLogger(config)}, nil
		},
	})
}

type MySQLClient struct {
	config       config.ModuleConfig
	Dsn          string
	ctx          context.Context
	inputHandler common.InputHandler
	db           *sql.DB
	logger       *slog.Logger
	dbMu         sync.Mutex
	cancel       context.CancelFunc
}

func (mc *MySQLClient) Id() string {
	return mc.config.Id
}

func (mc *MySQLClient) Type() string {
	return mc.config.Type
}

func (mc *MySQLClient) Start(ctx context.Context, inputHandler common.InputHandler) error {
	mc.logger.Debug("running")
	mc.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	mc.ctx = moduleContext
	mc.cancel = cancel

	db, err := sql.Open("mysql", mc.Dsn)
	if err != nil {
		return fmt.Errorf("mysql.client error connecting to database: %w", err)
	}

	// TODO(jwetzell): make configurable
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	mc.dbMu.Lock()
	mc.db = db
	mc.dbMu.Unlock()
	<-mc.ctx.Done()
	mc.logger.Debug("done")
	return nil
}

func (mc *MySQLClient) Stop() {
	if mc.cancel != nil {
		defer mc.cancel()
	}
	mc.dbMu.Lock()
	defer mc.dbMu.Unlock()
	if mc.db != nil {
		mc.db.Close()
	}
}

func (mc *MySQLClient) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	mc.dbMu.Lock()
	defer mc.dbMu.Unlock()
	if mc.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return mc.db.QueryContext(ctx, query, args...)
}
