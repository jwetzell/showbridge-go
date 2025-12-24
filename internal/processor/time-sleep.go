package processor

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type MetaDelay struct {
	config   config.ProcessorConfig
	logger   *slog.Logger
	Duration time.Duration
}

func (md *MetaDelay) Process(ctx context.Context, payload any) (any, error) {
	time.Sleep(md.Duration)
	return payload, nil
}

func (md *MetaDelay) Type() string {
	return md.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "time.sleep",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			duration, ok := params["duration"]
			if !ok {
				return nil, errors.New("time.sleep requires a duration parameter")
			}

			durationNum, ok := duration.(float64)

			if !ok {
				return nil, errors.New("time.sleep duration must be number")
			}

			return &MetaDelay{config: config, Duration: time.Millisecond * time.Duration(durationNum), logger: slog.Default().With("component", "processor")}, nil
		},
	})
}
