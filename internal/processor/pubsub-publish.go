package processor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"text/template"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type PubSubPublish struct {
	config   config.ProcessorConfig
	ModuleId string
	Topic    *template.Template
	logger   *slog.Logger
	module   common.PubSubModule
}

func (psp *PubSubPublish) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	if psp.module == nil {
		if wrappedPayload.Modules == nil {
			wrappedPayload.End = true
			return wrappedPayload, errors.New("pubsub.publish wrapped payload has no modules")
		}

		module, ok := wrappedPayload.Modules[psp.ModuleId]
		if !ok {
			wrappedPayload.End = true
			return wrappedPayload, fmt.Errorf("pubsub.publish unable to find module with id: %s", psp.ModuleId)
		}

		dbModule, ok := module.(common.PubSubModule)
		if !ok {
			wrappedPayload.End = true
			return wrappedPayload, fmt.Errorf("pubsub.publish module with id %s is not an OutputModule", psp.ModuleId)
		}
		psp.module = dbModule
	}

	var topicBuffer bytes.Buffer
	err := psp.Topic.Execute(&topicBuffer, wrappedPayload)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	err = psp.module.Publish(ctx, topicBuffer.String(), wrappedPayload.Payload)
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("pubsub.publish error publishing: %w", err)
	}

	return wrappedPayload, nil
}

func (psp *PubSubPublish) Type() string {
	return psp.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "pubsub.publish",
		Title: "Publish to Pub/Sub Topic",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"module": {
					Title:       "Module ID",
					Type:        "string",
					Description: "ID of the module to publish to",
				},
				"topic": {
					Title:       "Topic",
					Type:        "string",
					Description: "Topic to publish to",
				},
			},
			Required:             []string{"module", "topic"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ProcessorConfig) (Processor, error) {

			params := config.Params

			moduleIdString, err := params.GetString("module")
			if err != nil {
				return nil, fmt.Errorf("pubsub.publish module error: %w", err)
			}

			topicString, err := params.GetString("topic")
			if err != nil {
				return nil, fmt.Errorf("pubsub.publish topic error: %w", err)
			}

			topicTemplate, err := template.New("topic").Parse(topicString)

			if err != nil {
				return nil, err
			}
			return &PubSubPublish{config: config, ModuleId: moduleIdString, Topic: topicTemplate, logger: slog.Default().With("component", "processor", "type", config.Type)}, nil
		},
	})
}
