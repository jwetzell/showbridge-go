package routing

import (
	"context"
	"database/sql"

	"gitlab.com/gomidi/midi/v2"
)

type Module interface {
	Start(context.Context, Router) error
	Stop()
}

type ModuleManager interface {
	GetModule(string) (Module, bool)
}

type WriterModule interface {
	Write(context.Context, []byte) error
}

type MIDIOutputModule interface {
	Output(context.Context, midi.Message) error
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
