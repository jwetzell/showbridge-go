package processor

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "debug.log",
		Title: "Debug Log",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &DebugLog{config: config, logger: slog.Default().With("component", "processor", "type", config.Type)}, nil
		},
	})
}

type DebugLog struct {
	config config.ProcessorConfig
	logger *slog.Logger
}

// TODO(jwetzell): maybe make a more useful logging processor
func (dl *DebugLog) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadString := fmt.Sprintf("%+v", payload)
	payloadType := fmt.Sprintf("%T", payload)
	dl.logger.Debug("", "payload", payloadString, "payloadType", payloadType)
	return wrappedPayload, nil
}

func (dl *DebugLog) Type() string {
	return dl.config.Type
}

