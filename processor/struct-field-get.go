package processor

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:        "struct.field.get",
		Title:       "Get Struct Field",
		Description: "Extracts a field from a struct",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"name": {
					Title:       "Field Name",
					Description: "name of the struct field to extract",
					Type:        "string",
				},
			},
			Required:             []string{"name"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			processorParams := config.Params
			nameString, err := processorParams.GetString("name")
			if err != nil {
				return nil, fmt.Errorf("struct.field.get name error: %w", err)
			}

			return &StructFieldGet{config: config, Name: nameString}, nil
		},
	})
}

type StructFieldGet struct {
	config config.ProcessorConfig
	Name   string
}

func (sfg *StructFieldGet) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	s := reflect.ValueOf(payload)

	if s.Kind() != reflect.Struct {
		if s.Kind() == reflect.Pointer && s.Elem().Kind() == reflect.Struct {
			s = s.Elem()
		} else {
			wrappedPayload.End = true
			return wrappedPayload, errors.New("struct.field.get processor only accepts a struct payload")
		}
	}

	field := s.FieldByName(sfg.Name)
	if !field.IsValid() {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("struct.field.get field '%s' does not exist", sfg.Name)
	}

	wrappedPayload.Payload = field.Interface()
	return wrappedPayload, nil
}
