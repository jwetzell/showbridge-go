package processor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"text/template"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/common"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "os.exec",
		Title: "Exec Command",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"command": {
					Title:       "Command",
					Description: "the command to execute",
					Type:        "string",
				},
				"args": {
					Title:       "Args",
					Description: "the arguments for the command",
					Type:        "array",
					Items: &jsonschema.Schema{
						Type: "string",
					},
				},
			},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			Required:             []string{"command"},
		},
		New: func(moduleConfig config.ProcessorConfig) (Processor, error) {
			params := moduleConfig.Params

			commandString, err := params.GetString("command")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					return nil, fmt.Errorf("os.exec command error: not found")
				} else {
					return nil, fmt.Errorf("os.exec command error: %w", err)
				}
			}
			commandTemplate, err := template.New("command").Parse(commandString)

			if err != nil {
				return nil, err
			}

			argStrings, err := params.GetStringSlice("args")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					argStrings = []string{}
				} else {
					return nil, fmt.Errorf("os.exec args error: %w", err)
				}
			}

			argTemplates := make([]*template.Template, len(argStrings))
			for i, argString := range argStrings {
				argTemplate, err := template.New(fmt.Sprintf("arg-%d", i)).Parse(argString)
				if err != nil {
					return nil, err
				}
				argTemplates[i] = argTemplate
			}
			return &OsExec{config: moduleConfig, CommandTemplate: commandTemplate, Args: argTemplates}, nil
		},
	})
}

type OsExec struct {
	CommandTemplate *template.Template
	config          config.ProcessorConfig
	Args            []*template.Template
}

func (oe *OsExec) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	templateData := wrappedPayload

	var commandBuffer bytes.Buffer
	err := oe.CommandTemplate.Execute(&commandBuffer, templateData)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	commandString := commandBuffer.String()

	argStrings := make([]string, len(oe.Args))
	for i, argTemplate := range oe.Args {
		var argBuffer bytes.Buffer
		err := argTemplate.Execute(&argBuffer, templateData)
		if err != nil {
			wrappedPayload.End = true
			return wrappedPayload, err
		}
		argStrings[i] = argBuffer.String()
	}

	out, err := exec.Command(commandString, argStrings...).Output()
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("os.exec command execution error: %w", err)
	}
	wrappedPayload.Payload = out
	return wrappedPayload, nil
}

func (oe *OsExec) Id() string {
	return oe.config.Id
}

func (oe *OsExec) Type() string {
	return oe.config.Type
}
