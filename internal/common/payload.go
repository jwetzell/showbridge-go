package common

import (
	"context"
)

type WrappedPayload struct {
	Payload any
	Modules map[string]Module
	Sender  any
	Source  string
	End     bool
}

func GetWrappedPayload(ctx context.Context, payload any) WrappedPayload {
	wrappedPayload := WrappedPayload{
		Payload: payload,
		End:     false,
	}
	modules := ctx.Value(ModulesContextKey)
	if modules != nil {
		moduleMap, ok := modules.(map[string]Module)
		if ok {
			wrappedPayload.Modules = moduleMap
		} else {
			wrappedPayload.Modules = make(map[string]Module)
		}
	}

	sender := ctx.Value(SenderContextKey)
	if sender != nil {
		wrappedPayload.Sender = sender
	}

	source := ctx.Value(SourceContextKey)
	if source != nil {
		wrappedPayload.Source = source.(string)
	}
	return wrappedPayload
}
