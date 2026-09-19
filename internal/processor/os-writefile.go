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
	"github.com/jwetzell/showbridge-go/internal/common"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "os.writefile",
		Title: "Write File",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"path": {
					Title:       "Path",
					Description: "the path of the file to write to",
					Type:        "string",
				},
			},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			Required:             []string{"path"},
		},
		New: func(moduleConfig config.ProcessorConfig) (Processor, error) {
			params := moduleConfig.Params

			pathString, err := params.GetString("path")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					return nil, fmt.Errorf("os.writefile path error: not found")
				} else {
					return nil, fmt.Errorf("os.writefile path error: %w", err)
				}
			}
			pathTemplate, err := template.New("path").Parse(pathString)

			if err != nil {
				return nil, err
			}
			return &OsWriteFile{config: moduleConfig, PathTemplate: pathTemplate}, nil
		},
	})
}

type OsWriteFile struct {
	PathTemplate *template.Template
	config       config.ProcessorConfig
}

func (owf *OsWriteFile) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	templateData := wrappedPayload

	var pathBuffer bytes.Buffer
	err := owf.PathTemplate.Execute(&pathBuffer, templateData)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	pathString := pathBuffer.String()

	payloadBytes, ok := common.GetAnyAsByteSlice(wrappedPayload.Payload)

	if !ok {
		payloadString, ok := common.GetAnyAs[string](wrappedPayload.Payload)
		if !ok {
			wrappedPayload.End = true
			return wrappedPayload, errors.New("os.writefile can only process a string or []byte")
		}
		payloadBytes = []byte(payloadString)
	}

	// TODO(jwetzell): make perm configurable
	err = os.WriteFile(pathString, payloadBytes, 0644)
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("os.writefile error: %w", err)
	}
	return wrappedPayload, nil
}

func (owf *OsWriteFile) Id() string {
	return owf.config.Id
}

func (owf *OsWriteFile) Type() string {
	return owf.config.Type
}
