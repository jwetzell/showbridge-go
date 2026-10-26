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
		Type:        "struct.method.get",
		Title:       "Get Struct Method",
		Description: "Extracts and calls a method from a struct with no arguments",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"name": {
					Title:       "Method Name",
					Description: "name of the struct method to extract and call (with no arguments)",
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
				return nil, fmt.Errorf("struct.method.get name error: %w", err)
			}

			return &StructMethodGet{config: config, Name: nameString}, nil
		},
	})
}

type StructMethodGet struct {
	config config.ProcessorConfig
	Name   string
}

func (smg *StructMethodGet) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
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

	method := s.MethodByName(smg.Name)
	if !method.IsValid() {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("struct.method.get method '%s' does not exist", smg.Name)
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
