package processor

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/jwetzell/showbridge-go/internal/common"
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

func (srac *SipResponseAudioCreate) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {

	templateData := wrappedPayload

	var audioFileBuffer bytes.Buffer
	err := srac.AudioFile.Execute(&audioFileBuffer, templateData)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	audioFileString := audioFileBuffer.String()

	wrappedPayload.Payload = SipAudioFileResponse{
		PreWait:   srac.PreWait,
		PostWait:  srac.PostWait,
		AudioFile: audioFileString,
	}
	return wrappedPayload, nil
}

func (srac *SipResponseAudioCreate) Type() string {
	return srac.config.Type
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
