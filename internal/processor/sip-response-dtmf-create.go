package processor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"regexp"
	"text/template"

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

func (scc *SipResponseDTMFCreate) Process(ctx context.Context, payload any) (any, error) {

	var digitsBuffer bytes.Buffer
	err := scc.Digits.Execute(&digitsBuffer, payload)

	if err != nil {
		return nil, err
	}

	digitsString := digitsBuffer.String()

	if !scc.validDTMF.MatchString(digitsString) {
		return nil, errors.New("sip.response.dtmf.create result of digits template contains invalid characters")
	}

	return SipDTMFResponse{
		PreWait:  scc.PreWait,
		PostWait: scc.PostWait,
		Digits:   digitsString,
	}, nil
}

func (scc *SipResponseDTMFCreate) Type() string {
	return scc.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "sip.response.dtmf.create",
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
