package processor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"text/template"
	"time"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type HTTPRequestDo struct {
	config config.ProcessorConfig
	client *http.Client
	Method string
	URL    *template.Template
}

func (hrd *HTTPRequestDo) Process(ctx context.Context, payload any) (any, error) {

	templateData := GetTemplateData(ctx, payload)

	var urlBuffer bytes.Buffer
	err := hrd.URL.Execute(&urlBuffer, templateData)

	if err != nil {
		return nil, err
	}

	urlString := urlBuffer.String()

	//TODO(jwetzell): support body
	request, err := http.NewRequest(hrd.Method, urlString, bytes.NewBuffer([]byte{}))

	if err != nil {
		return nil, err
	}

	response, err := hrd.client.Do(request)

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	//TODO(jwetzell): support headers, etc
	return HTTPResponse{
		Status: response.StatusCode,
		Body:   body,
	}, nil
}

func (hrd *HTTPRequestDo) Type() string {
	return hrd.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "http.request.do",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			methodString, err := params.GetString("method")
			if err != nil {
				return nil, fmt.Errorf("http.request.do method error: %w", err)
			}

			urlString, err := params.GetString("url")
			if err != nil {
				return nil, fmt.Errorf("http.request.do url error: %w", err)
			}

			urlTemplate, err := template.New("url").Parse(urlString)

			if err != nil {
				return nil, err
			}
			client := &http.Client{
				Timeout: 10 * time.Second,
			}

			return &HTTPRequestDo{config: config, URL: urlTemplate, Method: methodString, client: client}, nil
		},
	})
}
