package module

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

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
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "db.sqlite",
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

func (t *DbSqlite) Start(ctx context.Context) error {
	t.logger.Debug("running")
	router, ok := ctx.Value(common.RouterContextKey).(common.RouteIO)

	if !ok {
		return errors.New("db.sqlite unable to get router from context")
	}
	t.router = router
	t.ctx = ctx

	db, err := sql.Open("sqlite", t.Dsn)
	if err != nil {
		return fmt.Errorf("db.sqlite error opening database: %w", err)
	}
	t.db = db
	defer t.db.Close()
	<-ctx.Done()
	return nil
}

func (t *DbSqlite) Stop() {
	if t.db != nil {
		t.db.Close()
	}
}

func (t *DbSqlite) Database() *sql.DB {
	return t.db
}
