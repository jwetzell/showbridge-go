//go:build js

package processor

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"syscall/js"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type WebSet struct {
	config    config.ProcessorConfig
	ModuleId  string
	ElementId string
	Property  string
	Value     *template.Template
	logger    *slog.Logger
}

func (kvs *WebSet) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {

	element := js.Global().Get("document").Call("getElementById", kvs.ElementId)

	if element.IsNull() || element.IsUndefined() {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("web.set unable to find element with id: %s", kvs.ElementId)
	}

	var valueBuffer bytes.Buffer
	err := kvs.Value.Execute(&valueBuffer, wrappedPayload)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	element.Set(kvs.Property, valueBuffer.String())
	return wrappedPayload, nil
}

func (kvs *WebSet) Type() string {
	return kvs.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "web.set",
		Title: "Set Web Element Property",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"id": {
					Title: "Element ID",
					Type:  "string",
				},
				"property": {
					Title: "Property",
					Type:  "string",
				},
				"value": {
					Title: "Value",
					Type:  "string",
				},
			},
			Required:             []string{"id", "property", "value"},
			AdditionalProperties: nil,
		},
		New: func(config config.ProcessorConfig) (Processor, error) {

			params := config.Params

			idString, err := params.GetString("id")
			if err != nil {
				return nil, fmt.Errorf("web.set id error: %w", err)
			}

			propertyString, err := params.GetString("property")
			if err != nil {
				return nil, fmt.Errorf("web.set property error: %w", err)
			}

			valueString, err := params.GetString("value")
			if err != nil {
				return nil, fmt.Errorf("web.set value error: %w", err)
			}
			valueTemplate, err := template.New("template").Parse(valueString)

			if err != nil {
				return nil, err
			}

			return &WebSet{config: config, ElementId: idString, Property: propertyString, Value: valueTemplate, logger: slog.Default().With("component", "processor", "type", config.Type)}, nil
		},
	})
}
