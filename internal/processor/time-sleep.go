package processor

import (
	"context"
	"fmt"
	"time"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type TimeSleep struct {
	config   config.ProcessorConfig
	Duration time.Duration
}

func (ts *TimeSleep) Process(ctx context.Context, payload any) (any, error) {
	time.Sleep(ts.Duration)
	return payload, nil
}

func (ts *TimeSleep) Type() string {
	return ts.config.Type
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

			return &TimeSleep{config: config, Duration: time.Millisecond * time.Duration(durationNum)}, nil
		},
	})
}
