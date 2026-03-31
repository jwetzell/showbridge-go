package processor

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type StructFieldGet struct {
	config config.ProcessorConfig
	Name   string
}

func (sf *StructFieldGet) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
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

	field := s.FieldByName(sf.Name)
	if !field.IsValid() {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("struct.field.get field '%s' does not exist", sf.Name)
	}

	wrappedPayload.Payload = field.Interface()
	return wrappedPayload, nil
}

func (sf *StructFieldGet) Type() string {
	return sf.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "struct.field.get",
		Title: "Get Struct Field",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"name": {
					Title: "Field Name",
					Type:  "string",
				},
			},
			Required:             []string{"name"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params
			nameString, err := params.GetString("name")
			if err != nil {
				return nil, fmt.Errorf("struct.field.get name error: %w", err)
			}

			return &StructFieldGet{config: config, Name: nameString}, nil
		},
	})
}
