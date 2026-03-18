package common

import (
	"context"
)

type Module interface {
	Id() string
	Type() string
	Start(context.Context) error
	Stop()
	Output(context.Context, any) error
}

type KeyValueModule interface {
	Get(key string) (any, error)
	Set(key string, value any) error
}
