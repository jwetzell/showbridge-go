package processor

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type HTTPResponseEncode struct {
	config config.ProcessorConfig
}

func (hre *HTTPResponseEncode) Process(ctx context.Context, payload any) (any, error) {
	payloadResponse, ok := payload.(*http.Response)

	if !ok {
		return nil, errors.New("http.response.encode processor only accepts an http.Response")
	}
	defer payloadResponse.Body.Close()

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
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &HTTPResponseEncode{config: config}, nil
		},
	})
}
