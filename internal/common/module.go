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
