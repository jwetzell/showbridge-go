package processor

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type DebugLog struct {
	config config.ProcessorConfig
	logger *slog.Logger
}

func (dl *DebugLog) Process(ctx context.Context, payload any) (any, error) {
	dl.logger.Debug("", "payload", payload, "payloadType", fmt.Sprintf("%T", payload))
	return payload, nil
}

func (dl *DebugLog) Type() string {
	return dl.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "debug.log",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &DebugLog{config: config, logger: slog.Default().With("component", "processor", "type", config.Type)}, nil
		},
	})
}
