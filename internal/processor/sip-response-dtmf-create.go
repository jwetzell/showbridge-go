//go:build !js

package processor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"regexp"
	"text/template"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type SipResponseDTMFCreate struct {
	config    config.ProcessorConfig
	PreWait   int
	PostWait  int
	Digits    *template.Template
	validDTMF *regexp.Regexp
}

type SipDTMFResponse struct {
	PreWait  int
	PostWait int
	Digits   string
}

func (srdc *SipResponseDTMFCreate) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {

	templateData := wrappedPayload

	var digitsBuffer bytes.Buffer
	err := srdc.Digits.Execute(&digitsBuffer, templateData)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	digitsString := digitsBuffer.String()

	if !srdc.validDTMF.MatchString(digitsString) {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("sip.response.dtmf.create result of digits template contains invalid characters")
	}

	wrappedPayload.Payload = SipDTMFResponse{
		PreWait:  srdc.PreWait,
		PostWait: srdc.PostWait,
		Digits:   digitsString,
	}
	return wrappedPayload, nil
}

func (srdc *SipResponseDTMFCreate) Type() string {
	return srdc.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "sip.response.dtmf.create",
		Title: "Create SIP DTMF Response",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"preWait": {
					Title: "Pre Wait (ms)",
					Type:  "integer",
				},
				"postWait": {
					Title: "Post Wait (ms)",
					Type:  "integer",
				},
				"digits": {
					Type: "string",
				},
			},
			Required:             []string{"preWait", "postWait", "digits"},
			AdditionalProperties: nil,
		},
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			preWaitNum, err := params.GetInt("preWait")
			if err != nil {
				return nil, fmt.Errorf("sip.response.dtmf.create preWait error: %w", err)
			}

			postWaitNum, err := params.GetInt("postWait")
			if err != nil {
				return nil, fmt.Errorf("sip.response.dtmf.create postWait error: %w", err)
			}

			digitsString, err := params.GetString("digits")
			if err != nil {
				return nil, fmt.Errorf("sip.response.dtmf.create digits error: %w", err)
			}

			digitsTemplate, err := template.New("digits").Parse(digitsString)

			if err != nil {
				return nil, err
			}
			return &SipResponseDTMFCreate{config: config, Digits: digitsTemplate, PreWait: int(preWaitNum), PostWait: int(postWaitNum), validDTMF: regexp.MustCompile(`^[0-9*#A-Da-d]+$`)}, nil
		},
	})
}
