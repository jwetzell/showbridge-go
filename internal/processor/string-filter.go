package processor

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type StringFilter struct {
	config  config.ProcessorConfig
	Pattern *regexp.Regexp
}

func (sf *StringFilter) Process(ctx context.Context, payload any) (any, error) {
	payloadString, ok := payload.(string)

	if !ok {
		return nil, errors.New("string.filter processor only accepts a string")
	}

	if !sf.Pattern.MatchString(payloadString) {
		return nil, nil
	}

	return payloadString, nil
}

func (sf *StringFilter) Type() string {
	return sf.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "string.filter",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			patternString, err := params.GetString("pattern")
			if err != nil {
				return nil, fmt.Errorf("string.filter pattern error: %w", err)
			}

			patternRegexp, err := regexp.Compile(patternString)

			if err != nil {
				return nil, err
			}

			return &StringFilter{config: config, Pattern: patternRegexp}, nil
		},
	})
}
