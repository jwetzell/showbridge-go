//go:build cgo

package processing

import (
	"context"
	"fmt"

	"gitlab.com/gomidi/midi/v2"
)

type MIDIMessageEncode struct {
	config ProcessorConfig
}

func (mme *MIDIMessageEncode) Process(ctx context.Context, payload any) (any, error) {
	payloadMessage, ok := payload.(midi.Message)

	if !ok {
		return nil, fmt.Errorf("midi.message.encode processor only accepts an midi.Message")
	}

	return payloadMessage.Bytes(), nil
}

func (mme *MIDIMessageEncode) Type() string {
	return mme.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "midi.message.encode",
		New: func(config ProcessorConfig) (Processor, error) {
			return &MIDIMessageEncode{config: config}, nil
		},
	})
}
