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
	"github.com/jwetzell/showbridge-go/params"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:        "os.exec",
		Title:       "Exec Command",
		Description: "Executes a specified command with optional arguments",
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
		New: func(processorConfig config.ProcessorConfig) (routing.Processor, error) {
			processorParams := processorConfig.Params

			commandString, err := processorParams.GetString("command")
			if err != nil {
				if errors.Is(err, params.ErrParamNotFound) {
					return nil, fmt.Errorf("os.exec command error: not found")
				} else {
					return nil, fmt.Errorf("os.exec command error: %w", err)
				}
			}
			commandTemplate, err := template.New("command").Parse(commandString)

			if err != nil {
				return nil, err
			}

			argStrings, err := processorParams.GetStringSlice("args")
			if err != nil {
				if errors.Is(err, params.ErrParamNotFound) {
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
			return &OsExec{config: processorConfig, CommandTemplate: commandTemplate, Args: argTemplates}, nil
		},
	})
}

type OsExec struct {
	CommandTemplate *template.Template
	config          config.ProcessorConfig
	Args            []*template.Template
}

func (oe *OsExec) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
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
