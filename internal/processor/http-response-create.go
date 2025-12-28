package processor

import (
	"bytes"
	"context"
	"errors"
	"text/template"

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

func (hre *HTTPResponseCreate) Process(ctx context.Context, payload any) (any, error) {
	var bodyBuffer bytes.Buffer
	err := hre.BodyTmpl.Execute(&bodyBuffer, payload)

	if err != nil {
		return nil, err
	}

	return HTTPResponse{
		Status: hre.Status,
		Body:   bodyBuffer.Bytes(),
	}, nil
}

func (hre *HTTPResponseCreate) Type() string {
	return hre.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "http.response.create",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			status, ok := params["status"]

			if !ok {
				return nil, errors.New("http.response.create requires a status parameter")
			}

			statusNum, ok := status.(float64)

			if !ok {
				return nil, errors.New("http.response.create status must be a number")
			}

			bodyTmpl, ok := params["bodyTemplate"]

			if !ok {
				return nil, errors.New("http.response.create requires a bodyTemplate parameter")
			}

			bodyTemplateString, ok := bodyTmpl.(string)

			if !ok {
				return nil, errors.New("http.response.create bodyTemplate must be a string")
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
