package processor

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type HTTPRequestFilter struct {
	config config.ProcessorConfig
	Path   *regexp.Regexp
	Method string
}

func (hrf *HTTPRequestFilter) Process(ctx context.Context, payload any) (any, error) {

	payloadRequest, ok := payload.(*http.Request)

	if !ok {
		return nil, errors.New("http.request.filter can only operate on http.Request payloads")
	}

	if hrf.Method != "" {
		if payloadRequest.Method != hrf.Method {
			return nil, nil
		}
	}

	if !hrf.Path.MatchString(payloadRequest.URL.Path) {
		return nil, nil
	}

	return payloadRequest, nil
}

func (hrf *HTTPRequestFilter) Type() string {
	return hrf.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "http.request.filter",
		New: func(moduleConfig config.ProcessorConfig) (Processor, error) {
			params := moduleConfig.Params
			pathString, err := params.GetString("path")
			if err != nil {
				return nil, fmt.Errorf("http.request.filter path error: %w", err)
			}

			pathRegexp, err := regexp.Compile(fmt.Sprintf("^%s$", pathString))

			if err != nil {
				return nil, err
			}

			methodString, err := params.GetString("method")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					return &HTTPRequestFilter{config: moduleConfig, Path: pathRegexp}, nil
				} else {
					return nil, fmt.Errorf("http.request.filter method error: %w", err)
				}
			}

			return &HTTPRequestFilter{config: moduleConfig, Path: pathRegexp, Method: methodString}, nil

		},
	})
}
