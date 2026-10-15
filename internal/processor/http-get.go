package processor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"text/template"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/common"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:        "http.get",
		Title:       "HTTP GET Request",
		Description: "Performs an HTTP GET request to the specified URL",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"url": {
					Title:       "URL",
					Description: "URL to send the HTTP request to",
					Type:        "string",
				},
				"headers": {
					Title:       "Headers",
					Description: "HTTP headers to include in the request",
					Type:        "object",
					AdditionalProperties: &jsonschema.Schema{
						Type: "string",
					},
					Default: json.RawMessage("{}"),
				},
			},
			Required:             []string{"url"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(processorConfig config.ProcessorConfig) (Processor, error) {
			params := processorConfig.Params

			urlString, err := params.GetString("url")
			if err != nil {
				return nil, fmt.Errorf("http.get url error: %w", err)
			}

			urlTemplate, err := template.New("url").Parse(urlString)

			if err != nil {
				return nil, err
			}
			client := &http.Client{
				Timeout: 10 * time.Second,
			}

			headers, err := params.GetStringMap("headers")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					headers = map[string]string{}
				} else {
					return nil, fmt.Errorf("http.get headers error: %w", err)
				}
			}
			return &HTTPGet{config: processorConfig, URL: urlTemplate, client: client, Headers: headers}, nil
		},
	})
}

type HTTPGet struct {
	config  config.ProcessorConfig
	client  *http.Client
	URL     *template.Template
	Headers map[string]string
}

func (hrd *HTTPGet) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {

	templateData := wrappedPayload

	var urlBuffer bytes.Buffer
	err := hrd.URL.Execute(&urlBuffer, templateData)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	urlString := urlBuffer.String()

	request, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}
	for key, value := range hrd.Headers {
		request.Header.Set(key, value)
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

	// TODO(jwetzell): support headers, etc
	wrappedPayload.Payload = httpResponse{
		Status:     response.Status,
		StatusCode: response.StatusCode,
		Proto:      response.Proto,
		Body:       body,
	}
	return wrappedPayload, nil
}

func (hrd *HTTPGet) Id() string {
	return hrd.config.Id
}

func (hrd *HTTPGet) Type() string {
	return hrd.config.Type
}
