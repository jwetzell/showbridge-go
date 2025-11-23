package processing

import (
	"context"
	"fmt"
	"regexp"
)

type StringFilter struct {
	config  ProcessorConfig
	Pattern *regexp.Regexp
}

func (se *StringFilter) Process(ctx context.Context, payload any) (any, error) {
	payloadString, ok := payload.(string)

	if !ok {
		return nil, fmt.Errorf("string.filter processor only accepts a string")
	}

	if !se.Pattern.MatchString(payloadString) {
		return nil, nil
	}

	return payloadString, nil
}

func (se *StringFilter) Type() string {
	return se.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "string.filter",
		New: func(config ProcessorConfig) (Processor, error) {
			params := config.Params

			pattern, ok := params["pattern"]

			if !ok {
				return nil, fmt.Errorf("http.request.filter requires an pattern parameter")
			}

			patternString, ok := pattern.(string)

			if !ok {
				return nil, fmt.Errorf("http.request.filter pattern must be a string")
			}

			patternRegexp, err := regexp.Compile(fmt.Sprintf("^%s$", patternString))

			if err != nil {
				return nil, err
			}

			return &StringFilter{config: config, Pattern: patternRegexp}, nil
		},
	})
}
