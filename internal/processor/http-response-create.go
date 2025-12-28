package processor

import (
	"bytes"
	"context"
	"errors"
	"text/template"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type HTTPResponseCreate struct {
	Status int
	Body   *template.Template
	config config.ProcessorConfig
}

type HTTPResponse struct {
	Status int
	Body   []byte
}

func (hre *HTTPResponseCreate) Process(ctx context.Context, payload any) (any, error) {
	var bodyBuffer bytes.Buffer
	err := hre.Body.Execute(&bodyBuffer, payload)

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
				return nil, errors.New("http.resposne.create status must be a number")
			}

			body, ok := params["body"]

			if !ok {
				return nil, errors.New("osc.message.create requires an body parameter")
			}

			bodyString, ok := body.(string)

			if !ok {
				return nil, errors.New("osc.message.create body must be a string")
			}

			bodyTemplate, err := template.New("body").Parse(bodyString)

			if err != nil {
				return nil, err
			}

			// TODO(jwetzell): support other body kind (direct bytes from input, from file?)
			return &HTTPResponseCreate{config: config, Status: int(statusNum), Body: bodyTemplate}, nil
		},
	})
}
