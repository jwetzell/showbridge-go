package test

import (
	"context"

	"github.com/jwetzell/showbridge-go/internal/common"
)

func GetContextWithModules(ctx context.Context, modules map[string]common.Module) context.Context {
	ctx = context.WithValue(ctx, common.ModulesContextKey, modules)
	return ctx
}

func GetContextWithRouter(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, common.RouterContextKey, GetNewTestRouter())
	return ctx
}

func GetContextWithSender(ctx context.Context, sender any) context.Context {
	ctx = context.WithValue(ctx, common.SenderContextKey, sender)
	return ctx
}

func GetContextWithSource(ctx context.Context, source string) context.Context {
	ctx = context.WithValue(ctx, common.SourceContextKey, source)
	return ctx
}
