package processor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"text/template"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/common"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:        "http.response",
		Title:       "HTTP Response",
		Description: "Respond to the current HTTP request if there is one",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"status": {
					Title:       "Status Code",
					Description: "status code to set on the response",
					Type:        "integer",
				},
				"bodyTemplate": {
					Title:       "Body Template",
					Description: "template for the response body",
					Type:        "string",
				},
			},
			Required:             []string{"status", "bodyTemplate"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ProcessorConfig) (common.Processor, error) {
			params := config.Params

			statusNum, err := params.GetInt("status")
			if err != nil {
				return nil, fmt.Errorf("http.response status error: %w", err)
			}

			bodyTemplateString, err := params.GetString("bodyTemplate")
			if err != nil {
				return nil, fmt.Errorf("http.response bodyTemplate error: %w", err)
			}

			bodyTemplate, err := template.New("body").Parse(bodyTemplateString)

			if err != nil {
				return nil, err
			}

			// TODO(jwetzell): support other body kind (direct bytes from input, from file?)
			return &HTTPResponse{config: config, Status: int(statusNum), BodyTmpl: bodyTemplate}, nil
		},
	})
}

type HTTPResponse struct {
	Status   int
	BodyTmpl *template.Template
	config   config.ProcessorConfig
}

type httpResponse struct {
	Status     string
	StatusCode int
	Proto      string
	Body       []byte
}

func (hrc *HTTPResponse) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	templateData := wrappedPayload

	var bodyBuffer bytes.Buffer
	err := hrc.BodyTmpl.Execute(&bodyBuffer, templateData)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	responseWriter, ok := ctx.Value(common.HttpServerResponseContextKey).(*common.HTTPServerResponseWriter)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("http.response input must come from a http.server")
	}

	locked := responseWriter.TryLock()
	if !locked {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("http.response could not acquire lock on response writer")
	}
	defer responseWriter.Unlock()

	if responseWriter.Done() {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("http.response response writer has already been written to")
	}

	responseWriter.WriteHeader(hrc.Status)
	_, err = responseWriter.Write(bodyBuffer.Bytes())
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	// wrappedPayload.Payload = httpResponse{
	// 	Status: hrc.Status,
	// 	Body:   bodyBuffer.Bytes(),
	// }
	return wrappedPayload, nil
}
