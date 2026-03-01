package processor

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"text/template"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type HTTPRequestCreate struct {
	config config.ProcessorConfig
	Method string
	URL    *template.Template
}

func (hrc *HTTPRequestCreate) Process(ctx context.Context, payload any) (any, error) {

	var urlBuffer bytes.Buffer
	err := hrc.URL.Execute(&urlBuffer, payload)

	if err != nil {
		return nil, err
	}

	urlString := urlBuffer.String()

	//TODO(jwetzell): support body
	request, err := http.NewRequest(hrc.Method, urlString, bytes.NewBuffer([]byte{}))

	if err != nil {
		return nil, err
	}

	return request, nil
}

func (hrc *HTTPRequestCreate) Type() string {
	return hrc.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "http.request.create",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			methodString, err := params.GetString("method")
			if err != nil {
				return nil, fmt.Errorf("http.request.create method error: %w", err)
			}

			urlString, err := params.GetString("url")
			if err != nil {
				return nil, fmt.Errorf("http.request.create url error: %w", err)
			}

			urlTemplate, err := template.New("url").Parse(urlString)

			if err != nil {
				return nil, err
			}
			return &HTTPRequestCreate{config: config, URL: urlTemplate, Method: methodString}, nil
		},
	})
}
