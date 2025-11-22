package processing

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type HTTPRequestEncode struct {
	config ProcessorConfig
}

func (hre *HTTPRequestEncode) Process(ctx context.Context, payload any) (any, error) {
	payloadRequest, ok := payload.(*http.Request)

	if !ok {
		return nil, fmt.Errorf("http.request.encode processor only accepts an http.Request")
	}

	bytes, err := io.ReadAll(payloadRequest.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (hre *HTTPRequestEncode) Type() string {
	return hre.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "http.request.encode",
		New: func(config ProcessorConfig) (Processor, error) {
			return &HTTPRequestEncode{config: config}, nil
		},
	})
}
