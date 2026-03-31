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

type StructMethodGet struct {
	config config.ProcessorConfig
	Name   string
}

func (sm *StructMethodGet) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	s := reflect.ValueOf(payload)

	if s.Kind() != reflect.Struct {
		if s.Kind() == reflect.Pointer && s.Elem().Kind() == reflect.Struct {
			s = s.Elem()
		} else {
			wrappedPayload.End = true
			return wrappedPayload, errors.New("struct.method.get processor only accepts a struct payload")
		}
	}

	method := s.MethodByName(sm.Name)
	if !method.IsValid() {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("struct.method.get method '%s' does not exist", sm.Name)
	}

	value := method.Call(nil)

	if len(value) == 0 {
		wrappedPayload.End = true
		wrappedPayload.Payload = nil
		return wrappedPayload, nil
	}

	if len(value) == 1 {
		wrappedPayload.Payload = value[0].Interface()
		return wrappedPayload, nil
	}

	results := make([]any, len(value))

	for i, v := range value {
		results[i] = v.Interface()
	}

	wrappedPayload.Payload = results
	return wrappedPayload, nil
}

func (sm *StructMethodGet) Type() string {
	return sm.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "struct.method.get",
		Title: "Get Struct Method",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"name": {
					Title: "Method Name",
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
				return nil, fmt.Errorf("struct.method.get name error: %w", err)
			}

			return &StructMethodGet{config: config, Name: nameString}, nil
		},
	})
}
