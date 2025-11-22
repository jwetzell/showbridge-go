package processing

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type HTTPResponseEncode struct {
	config ProcessorConfig
}

func (hre *HTTPResponseEncode) Process(ctx context.Context, payload any) (any, error) {
	payloadResponse, ok := payload.(*http.Response)
	defer payloadResponse.Body.Close()

	if !ok {
		return nil, fmt.Errorf("http.response.encode processor only accepts an http.Response")
	}

	bytes, err := io.ReadAll(payloadResponse.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (hre *HTTPResponseEncode) Type() string {
	return hre.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "http.response.encode",
		New: func(config ProcessorConfig) (Processor, error) {
			return &HTTPResponseEncode{config: config}, nil
		},
	})
}
