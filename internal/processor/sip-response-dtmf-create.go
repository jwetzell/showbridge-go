package processor

import (
	"bytes"
	"context"
	"errors"
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

			preWait, ok := params["preWait"]

			if !ok {
				return nil, errors.New("sip.response.dtmf.create requires a preWait parameter")
			}

			preWaitNum, ok := preWait.(float64)

			if !ok {
				return nil, errors.New("sip.response.dtmf.create preWait must be a number")
			}

			postWait, ok := params["postWait"]

			if !ok {
				return nil, errors.New("sip.response.dtmf.create requires a postWait parameter")
			}

			postWaitNum, ok := postWait.(float64)

			if !ok {
				return nil, errors.New("sip.response.dtmf.create postWait must be a number")
			}

			digits, ok := params["digits"]

			if !ok {
				return nil, errors.New("sip.response.dtmf.create requires a digits parameter")
			}

			digitsString, ok := digits.(string)

			if !ok {
				return nil, errors.New("sip.response.dtmf.create digits must be a string")
			}

			digitsTemplate, err := template.New("digits").Parse(digitsString)

			if err != nil {
				return nil, err
			}
			return &SipResponseDTMFCreate{config: config, Digits: digitsTemplate, PreWait: int(preWaitNum), PostWait: int(postWaitNum), validDTMF: regexp.MustCompile(`^[0-9*#A-Da-d]+$`)}, nil
		},
	})
}
