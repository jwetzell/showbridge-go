package common

import (
	"context"
	"database/sql"
)

type Module interface {
	Id() string
	Type() string
	Start(context.Context) error
	Stop()
}

type OutputModule interface {
	Output(context.Context, any) error
}

type KeyValueModule interface {
	Get(key string) (any, error)
	Set(key string, value any) error
}

type DatabaseModule interface {
	Database() *sql.DB
}
