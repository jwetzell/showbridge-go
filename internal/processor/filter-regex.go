package processor

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type FilterRegex struct {
	config  config.ProcessorConfig
	Pattern *regexp.Regexp
}

func (fr *FilterRegex) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadString, ok := common.GetAnyAs[string](payload)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("filter.regex processor only accepts a string")
	}

	if !fr.Pattern.MatchString(payloadString) {
		wrappedPayload.End = true
		return wrappedPayload, nil
	}

	wrappedPayload.Payload = payloadString
	return wrappedPayload, nil
}

func (fr *FilterRegex) Type() string {
	return fr.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "filter.regex",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			patternString, err := params.GetString("pattern")
			if err != nil {
				return nil, fmt.Errorf("filter.regex pattern error: %w", err)
			}

			patternRegexp, err := regexp.Compile(patternString)

			if err != nil {
				return nil, err
			}

			return &FilterRegex{config: config, Pattern: patternRegexp}, nil
		},
	})
}
