package processor

import (
	"context"
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
		return nil, fmt.Errorf("http.request.filter can only operate on http.Request payloads")
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
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params
			path, ok := params["path"]

			if !ok {
				return nil, fmt.Errorf("http.request.filter requires an path parameter")
			}

			pathString, ok := path.(string)

			if !ok {
				return nil, fmt.Errorf("http.request.filter path must be a string")
			}

			pathRegexp, err := regexp.Compile(fmt.Sprintf("^%s$", pathString))

			if err != nil {
				return nil, err
			}

			method, ok := params["method"]

			if ok {
				methodString, ok := method.(string)

				if !ok {
					return nil, fmt.Errorf("http.request.filter method must be a string")
				}
				return &HTTPRequestFilter{config: config, Path: pathRegexp, Method: methodString}, nil
			}

			return &HTTPRequestFilter{config: config, Path: pathRegexp}, nil
		},
	})
}
