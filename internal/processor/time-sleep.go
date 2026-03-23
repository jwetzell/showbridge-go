package processor

import (
	"context"
	"fmt"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type TimeSleep struct {
	config   config.ProcessorConfig
	Duration time.Duration
}

func (ts *TimeSleep) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	time.Sleep(ts.Duration)
	return wrappedPayload, nil
}

func (ts *TimeSleep) Type() string {
	return ts.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "time.sleep",
		Title: "Sleep",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"duration": {
					Title:       "Duration",
					Type:        "integer",
					Description: "Duration to sleep in milliseconds",
				},
			},
			Required: []string{"duration"},
		},
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
