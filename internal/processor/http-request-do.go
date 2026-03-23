package processor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"text/template"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type HTTPRequestDo struct {
	config config.ProcessorConfig
	client *http.Client
	Method string
	URL    *template.Template
}

func (hrd *HTTPRequestDo) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {

	templateData := wrappedPayload

	var urlBuffer bytes.Buffer
	err := hrd.URL.Execute(&urlBuffer, templateData)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	urlString := urlBuffer.String()

	//TODO(jwetzell): support body
	request, err := http.NewRequest(hrd.Method, urlString, bytes.NewBuffer([]byte{}))

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	response, err := hrd.client.Do(request)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	//TODO(jwetzell): support headers, etc
	wrappedPayload.Payload = HTTPResponse{
		Status: response.StatusCode,
		Body:   body,
	}
	return wrappedPayload, nil
}

func (hrd *HTTPRequestDo) Type() string {
	return hrd.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "http.request.do",
		Title: "Do HTTP Request",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"method": {
					Title: "HTTP Method",
					Type:  "string",
					Enum:  []any{"GET", "POST"},
				},
				"url": {
					Title: "URL",
					Type:  "string",
				},
			},
			Required:             []string{"method", "url"},
			AdditionalProperties: nil,
		},
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
