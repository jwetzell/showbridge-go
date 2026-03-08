package processor

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"sync"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type Processor interface {
	Type() string
	Process(context.Context, any) (any, error)
}

type ProcessorRegistration struct {
	Type string `json:"type"`
	New  func(config.ProcessorConfig) (Processor, error)
}

func RegisterProcessor(processor ProcessorRegistration) {

	if processor.Type == "" {
		panic("processor type is missing")
	}
	if processor.New == nil {
		panic("missing ProcessorRegistration.New")
	}

	processorRegistryMu.Lock()
	defer processorRegistryMu.Unlock()

	if _, ok := ProcessorRegistry[string(processor.Type)]; ok {
		panic(fmt.Sprintf("processor already registered: %s", processor.Type))
	}
	ProcessorRegistry[string(processor.Type)] = processor
}

var (
	processorRegistryMu sync.RWMutex
	ProcessorRegistry   = make(map[string]ProcessorRegistration)
)

func GetAnyAs[T any](p any) (T, bool) {
	typed, ok := p.(T)
	return typed, ok
}

func GetAnyAsByteSlice(p any) ([]byte, bool) {
	v := reflect.ValueOf(p)
	if v.Kind() != reflect.Slice {
		return nil, false
	}

	result := make([]byte, v.Len())
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i).Interface()
		byteValue, ok := elem.(byte)
		if ok {
			result[i] = byteValue
			continue
		}
		uintValue, ok := elem.(uint)
		if ok {
			if uintValue > 255 {
				return nil, false
			}
			result[i] = byte(uintValue)
			continue
		}
		intValue, ok := elem.(int)
		if ok {
			if intValue < 0 || intValue > 255 {
				return nil, false
			}
			result[i] = byte(intValue)
			continue
		}
		floatValue, ok := elem.(float64)
		if ok {
			if floatValue != math.Floor(floatValue) {
				return nil, false
			}
			if floatValue < 0 || floatValue > 255 {
				return nil, false
			}
			result[i] = byte(floatValue)
			continue
		}
		return nil, false
	}
	return result, true
}

type TemplateData struct {
	Payload any
	Modules any
	Sender  any
}

func GetTemplateData(ctx context.Context, payload any) TemplateData {
	templateData := TemplateData{Payload: payload}
	modules := ctx.Value(common.ModulesContextKey)
	if modules != nil {
		templateData.Modules = modules
	}

	sender := ctx.Value(common.SenderContextKey)
	if sender != nil {
		templateData.Sender = sender
	}
	return templateData
}
