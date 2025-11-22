package showbridge

import (
	"fmt"

	"github.com/jwetzell/showbridge-go/internal/processing"
)

type RouteError struct {
	Index  int
	Config RouteConfig
	Error  error
}

type Route struct {
	index      int
	Input      string
	Processors []processing.Processor
	Output     string
	router     *Router
}

type RouteConfig struct {
	Input      string                       `json:"input"`
	Processors []processing.ProcessorConfig `json:"processors"`
	Output     string                       `json:"output"`
}

func NewRoute(index int, config RouteConfig, router *Router) (*Route, error) {
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

	return &Route{Input: config.Input, Processors: processors, Output: config.Output, router: router, index: index}, nil
}

func (r *Route) HandleInput(sourceId string, payload any) error {
	var err error
	for _, processor := range r.Processors {
		payload, err = processor.Process(r.router.Context, payload)
		if err != nil {
			return err
		}
		//NOTE(jwetzell) nil payload will result in the route being "terminated"
		if payload == nil {
			return nil
		}
	}
	return r.HandleOutput(payload)
}

func (r *Route) HandleOutput(payload any) error {
	return r.router.HandleOutput(r.Output, payload)
}
