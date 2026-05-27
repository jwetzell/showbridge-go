package common

import (
	"context"
	"database/sql"
)

type Module interface {
	Id() string
	Type() string
	Start(context.Context, InputHandler) error
	Stop()
}

type OutputModule interface {
	Output(context.Context, any) error
}

type KeyValueModule interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, value any) error
}

type DatabaseModule interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type PubSubModule interface {
	Publish(ctx context.Context, topic string, payload any) error
}
