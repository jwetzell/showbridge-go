package processing

import (
	"context"
	"fmt"
	"log/slog"
)

type DebugLog struct {
	config ProcessorConfig
}

func (dl *DebugLog) Process(ctx context.Context, payload any) (any, error) {
	slog.Debug("debug.log", "payload", payload, "payloadType", fmt.Sprintf("%T", payload))
	return payload, nil
}

func (dl *DebugLog) Type() string {
	return dl.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "debug.log",
		New: func(config ProcessorConfig) (Processor, error) {
			return &DebugLog{config: config}, nil
		},
	})
}
