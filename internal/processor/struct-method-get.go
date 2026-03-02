package processor

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type StructMethodGet struct {
	config config.ProcessorConfig
	Name   string
}

func (sm *StructMethodGet) Process(ctx context.Context, payload any) (any, error) {
	s := reflect.ValueOf(payload)

	if s.Kind() != reflect.Struct {
		return nil, errors.New("struct.method.get processor only accepts a struct payload")
	}

	method := s.MethodByName(sm.Name)
	if !method.IsValid() {
		return nil, fmt.Errorf("struct.method.get method '%s' does not exist", sm.Name)
	}

	value := method.Call(nil)

	if len(value) == 0 {
		return nil, nil
	}

	if len(value) == 1 {
		return value[0].Interface(), nil
	}

	results := make([]any, len(value))

	for i, v := range value {
		results[i] = v.Interface()
	}

	return results, nil
}

func (sm *StructMethodGet) Type() string {
	return sm.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "struct.method.get",
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
