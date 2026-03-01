package processor

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type SipResponseAudioCreate struct {
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

func (scc *SipResponseAudioCreate) Process(ctx context.Context, payload any) (any, error) {

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

func (scc *SipResponseAudioCreate) Type() string {
	return scc.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "sip.response.audio.create",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			preWaitNum, err := params.GetInt("preWait")
			if err != nil {
				return nil, fmt.Errorf("sip.response.audio.create preWait error: %w", err)
			}

			postWaitNum, err := params.GetInt("postWait")
			if err != nil {
				return nil, fmt.Errorf("sip.response.audio.create postWait error: %w", err)
			}

			audioFileString, err := params.GetString("audioFile")
			if err != nil {
				return nil, fmt.Errorf("sip.response.audio.create audioFile error: %w", err)
			}

			audioFileTemplate, err := template.New("audioFile").Parse(audioFileString)

			if err != nil {
				return nil, err
			}
			return &SipResponseAudioCreate{config: config, AudioFile: audioFileTemplate, PreWait: int(preWaitNum), PostWait: int(postWaitNum)}, nil
		},
	})
}
