package showbridge

import (
	"fmt"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processing"
)

type RouteError struct {
	Index  int
	Config config.RouteConfig
	Error  error
}

type Route interface {
	Input() string
	Output() string
	HandleInput(sourceId string, payload any, router *Router) error
	HandleOutput(sourceId string, payload any, router *Router) error
}

type ProcessorRoute struct {
	input      string
	processors []processing.Processor
	output     string
}

func NewRoute(config config.RouteConfig) (Route, error) {
	processors := []processing.Processor{}

	if len(config.Processors) > 0 {
		for _, processorDecl := range config.Processors {
			processorInfo, ok := processing.ProcessorRegistry[processorDecl.Type]
			if !ok {
				return nil, fmt.Errorf("problem loading processor registration for processor type: %s", processorDecl.Type)
			}

			processor, err := processorInfo.New(processorDecl)
			if err != nil {
				return nil, err
			}
			processors = append(processors, processor)
		}
	}

	return &ProcessorRoute{input: config.Input, processors: processors, output: config.Output}, nil
}

func (r *ProcessorRoute) Input() string {
	return r.input
}

func (r *ProcessorRoute) Output() string {
	return r.output
}

func (r *ProcessorRoute) HandleInput(sourceId string, payload any, router *Router) error {
	var err error
	for _, processor := range r.processors {
		payload, err = processor.Process(router.Context, payload)
		if err != nil {
			return err
		}
		//NOTE(jwetzell) nil payload will result in the route being "terminated"
		if payload == nil {
			return nil
		}
	}
	return r.HandleOutput(sourceId, payload, router)
}

func (r *ProcessorRoute) HandleOutput(sourceId string, payload any, router *Router) error {
	return router.HandleOutput(sourceId, r.output, payload)
}
