package processor

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:        "debug.log",
		Title:       "Debug Log",
		Description: "Logs the payload for debugging purposes",
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			return &DebugLog{config: config, logger: slog.Default().With("component", "processor", "type", config.Type)}, nil
		},
	})
}

type DebugLog struct {
	config config.ProcessorConfig
	logger *slog.Logger
}

// TODO(jwetzell): maybe make a more useful logging processor
func (dl *DebugLog) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadString := fmt.Sprintf("%+v", payload)
	payloadType := fmt.Sprintf("%T", payload)
	dl.logger.Debug("", "payload", payloadString, "payloadType", payloadType)
	return wrappedPayload, nil
}
