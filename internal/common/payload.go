package common

import (
	"context"
)

type WrappedPayload struct {
	Payload any
	Modules any
	Sender  any
	End     bool
}

func GetWrappedPayload(ctx context.Context, payload any) WrappedPayload {
	templateData := WrappedPayload{
		Payload: payload,
		End:     false,
	}
	modules := ctx.Value(ModulesContextKey)
	if modules != nil {
		templateData.Modules = modules
	}

	sender := ctx.Value(SenderContextKey)
	if sender != nil {
		templateData.Sender = sender
	}
	return templateData
}
