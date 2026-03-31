package processor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type RouterInput struct {
	config   config.ProcessorConfig
	SourceId string
	logger   *slog.Logger
}

func (ro *RouterInput) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {

	payload := wrappedPayload.Payload
	router, ok := ctx.Value(common.RouterContextKey).(common.RouteIO)
	if !ok {

		wrappedPayload.End = true
		return wrappedPayload, errors.New("router.input no router found")
	}

	_, err := router.HandleInput(ctx, ro.SourceId, payload)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("router.input failed to send input")
	}

	wrappedPayload.Payload = payload

	return wrappedPayload, nil
}

func (ro *RouterInput) Type() string {
	return ro.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "router.input",
		Title: "Router Input",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"source": {
					Title:       "Source",
					Type:        "string",
					Description: "source to report as to the router",
				},
			},
			Required:             []string{"source"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ProcessorConfig) (Processor, error) {

			params := config.Params

			sourceId, err := params.GetString("source")

			if err != nil {
				return nil, fmt.Errorf("router.input source error: %w", err)
			}

			return &RouterInput{config: config, SourceId: sourceId, logger: slog.Default().With("component", "processor", "type", config.Type)}, nil
		},
	})
}
