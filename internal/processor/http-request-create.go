package processor

import (
	"bytes"
	"context"
	"errors"
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

			method, ok := params["method"]

			if !ok {
				return nil, errors.New("http.request.create requires a method parameter")
			}

			methodString, ok := method.(string)

			if !ok {
				return nil, errors.New("http.request.create url must be a string")
			}

			url, ok := params["url"]

			if !ok {
				return nil, errors.New("http.request.create requires a url parameter")
			}

			urlString, ok := url.(string)

			if !ok {
				return nil, errors.New("http.request.create url must be a string")
			}

			urlTemplate, err := template.New("url").Parse(urlString)

			if err != nil {
				return nil, err
			}
			return &HTTPRequestCreate{config: config, URL: urlTemplate, Method: methodString}, nil
		},
	})
}
