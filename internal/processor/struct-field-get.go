package processor

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type StructFieldGet struct {
	config config.ProcessorConfig
	Name   string
}

func (sf *StructFieldGet) Process(ctx context.Context, payload any) (any, error) {
	s := reflect.ValueOf(payload)

	if s.Kind() != reflect.Struct {
		return nil, errors.New("struct.field.get processor only accepts a struct payload")
	}

	field := s.FieldByName(sf.Name)
	if !field.IsValid() {
		return nil, fmt.Errorf("struct.field.get field '%s' does not exist", sf.Name)
	}

	return field.Interface(), nil
}

func (sf *StructFieldGet) Type() string {
	return sf.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "struct.field.get",
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
