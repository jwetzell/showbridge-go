package processor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"text/template"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "pubsub.publish",
		Title: "Publish to Pub/Sub Topic",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"module": {
					Title:       "Module ID",
					Description: "ID of the module to publish to",
					Type:        "string",
				},
				"topic": {
					Title:       "Topic",
					Description: "topic to publish to",
					Type:        "string",
				},
			},
			Required:             []string{"module", "topic"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ProcessorConfig) (routing.Processor, error) {

			processorParams := config.Params

			moduleIdString, err := processorParams.GetString("module")
			if err != nil {
				return nil, fmt.Errorf("pubsub.publish module error: %w", err)
			}

			topicString, err := processorParams.GetString("topic")
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

type PubSubPublish struct {
	config   config.ProcessorConfig
	ModuleId string
	Topic    *template.Template
	logger   *slog.Logger
	module   routing.PubSubModule
}

func (psp *PubSubPublish) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, inputHandler routing.InputHandler) (routing.WrappedPayload, error) {
	if psp.module == nil {
		if moduleManager == nil {
			wrappedPayload.End = true
			return wrappedPayload, errors.New("pubsub.publish module manager is nil")
		}

		module, ok := moduleManager.GetModule(psp.ModuleId)
		if !ok {
			wrappedPayload.End = true
			return wrappedPayload, fmt.Errorf("pubsub.publish unable to find module with id: %s", psp.ModuleId)
		}

		pubSubModule, ok := module.(routing.PubSubModule)
		if !ok {
			wrappedPayload.End = true
			return wrappedPayload, fmt.Errorf("pubsub.publish module with id %s is not a PubSubModule", psp.ModuleId)
		}
		psp.module = pubSubModule
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
