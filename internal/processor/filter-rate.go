package processor

import (
	"context"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"golang.org/x/time/rate"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "filter.rate",
		Title: "Filter by Rate",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"rate": {
					Type:        "integer",
					Title:       "Rate",
					Description: "The number of events to allow per second.",
				},
			},
			Required: []string{"rate"},
		},
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params
			rateInt, err := params.GetInt("rate")
			if err != nil {
				return nil, fmt.Errorf("filter.rate rate error: %w", err)
			}

			limiter := rate.NewLimiter(rate.Limit(rateInt), rateInt*2)

			return &FilterRate{config: config, limiter: limiter}, nil
		},
	})
}

type FilterRate struct {
	config  config.ProcessorConfig
	limiter *rate.Limiter
}

func (fc *FilterRate) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	err := fc.limiter.Wait(ctx)
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}
	return wrappedPayload, nil
}

func (fc *FilterRate) Type() string {
	return fc.config.Type
}
