package processor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"text/template"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/params"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "os.readfile",
		Title: "Read File",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"path": {
					Title:       "Path",
					Description: "the path of the file to read",
					Type:        "string",
				},
			},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			Required:             []string{"path"},
		},
		New: func(processorConfig config.ProcessorConfig) (routing.Processor, error) {
			processorParams := processorConfig.Params

			pathString, err := processorParams.GetString("path")
			if err != nil {
				if errors.Is(err, params.ErrParamNotFound) {
					return nil, fmt.Errorf("os.readfile path error: not found")
				} else {
					return nil, fmt.Errorf("os.readfile path error: %w", err)
				}
			}
			pathTemplate, err := template.New("path").Parse(pathString)

			if err != nil {
				return nil, err
			}
			return &OsReadFile{config: processorConfig, PathTemplate: pathTemplate}, nil
		},
	})
}

type OsReadFile struct {
	PathTemplate *template.Template
	config       config.ProcessorConfig
}

func (orf *OsReadFile) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
	templateData := wrappedPayload

	var pathBuffer bytes.Buffer
	err := orf.PathTemplate.Execute(&pathBuffer, templateData)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	pathString := pathBuffer.String()

	data, err := os.ReadFile(pathString)
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("os.readfile error: %w", err)
	}
	wrappedPayload.Payload = data
	return wrappedPayload, nil
}
