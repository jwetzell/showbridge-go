package processor

import (
	"bytes"
	"context"
	"errors"
	"text/template"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type SipResponseCreate struct {
	config    config.ProcessorConfig
	PreWait   int
	PostWait  int
	AudioFile *template.Template
}

type SipAudioFileResponse struct {
	PreWait   int
	PostWait  int
	AudioFile string
}

func (scc *SipResponseCreate) Process(ctx context.Context, payload any) (any, error) {

	var audioFileBuffer bytes.Buffer
	err := scc.AudioFile.Execute(&audioFileBuffer, payload)

	if err != nil {
		return nil, err
	}

	audioFileString := audioFileBuffer.String()

	return SipAudioFileResponse{
		PreWait:   scc.PreWait,
		PostWait:  scc.PostWait,
		AudioFile: audioFileString,
	}, nil
}

func (scc *SipResponseCreate) Type() string {
	return scc.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "sip.response.create",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			preWait, ok := params["preWait"]

			if !ok {
				return nil, errors.New("sip.response.create requires a preWait parameter")
			}

			preWaitNum, ok := preWait.(float64)

			if !ok {
				return nil, errors.New("sip.response.create preWait must be a number")
			}

			postWait, ok := params["postWait"]

			if !ok {
				return nil, errors.New("sip.response.create requires a postWait parameter")
			}

			postWaitNum, ok := postWait.(float64)

			if !ok {
				return nil, errors.New("sip.response.create postWait must be a number")
			}

			audioFile, ok := params["audioFile"]

			if !ok {
				return nil, errors.New("sip.response.create requires a audioFile parameter")
			}

			audioFileString, ok := audioFile.(string)

			if !ok {
				return nil, errors.New("sip.response.create audioFile must be a string")
			}

			audioFileTemplate, err := template.New("audioFile").Parse(audioFileString)

			if err != nil {
				return nil, err
			}
			return &SipResponseCreate{config: config, AudioFile: audioFileTemplate, PreWait: int(preWaitNum), PostWait: int(postWaitNum)}, nil
		},
	})
}
