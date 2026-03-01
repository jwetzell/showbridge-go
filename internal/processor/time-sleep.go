package processor

import (
	"context"
	"fmt"
	"time"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type MetaDelay struct {
	config   config.ProcessorConfig
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

			durationNum, err := params.GetInt("duration")
			if err != nil {
				return nil, fmt.Errorf("time.sleep duration error: %w", err)
			}

			return &MetaDelay{config: config, Duration: time.Millisecond * time.Duration(durationNum)}, nil
		},
	})
}
