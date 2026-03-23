package processor

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type HTTPResponseCreate struct {
	Status   int
	BodyTmpl *template.Template
	config   config.ProcessorConfig
}

type HTTPResponse struct {
	Status int
	Body   []byte
}

func (hrc *HTTPResponseCreate) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	templateData := wrappedPayload

	var bodyBuffer bytes.Buffer
	err := hrc.BodyTmpl.Execute(&bodyBuffer, templateData)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}
	wrappedPayload.Payload = HTTPResponse{
		Status: hrc.Status,
		Body:   bodyBuffer.Bytes(),
	}
	return wrappedPayload, nil
}

func (hrc *HTTPResponseCreate) Type() string {
	return hrc.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "http.response.create",
		Title: "Create HTTP Response",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"status": {
					Title: "Status Code",
					Type:  "integer",
				},
				"body": {
					Title: "Body",
					Type:  "string",
				},
			},
			Required:             []string{"status", "body"},
			AdditionalProperties: nil,
		},
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			statusNum, err := params.GetInt("status")
			if err != nil {
				return nil, fmt.Errorf("http.response.create status error: %w", err)
			}

			bodyTemplateString, err := params.GetString("bodyTemplate")
			if err != nil {
				return nil, fmt.Errorf("http.response.create bodyTemplate error: %w", err)
			}

			bodyTemplate, err := template.New("body").Parse(bodyTemplateString)

			if err != nil {
				return nil, err
			}

			// TODO(jwetzell): support other body kind (direct bytes from input, from file?)
			return &HTTPResponseCreate{config: config, Status: int(statusNum), BodyTmpl: bodyTemplate}, nil
		},
	})
}
