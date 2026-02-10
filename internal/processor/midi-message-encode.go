//go:build cgo

package processor

import (
	"context"
	"errors"

	"github.com/jwetzell/showbridge-go/internal/config"
	"gitlab.com/gomidi/midi/v2"
)

type MIDIMessageEncode struct {
	config config.ProcessorConfig
}

func (mme *MIDIMessageEncode) Process(ctx context.Context, payload any) (any, error) {
	payloadMessage, ok := payload.(midi.Message)

	if !ok {
		return nil, errors.New("midi.message.encode processor only accepts a midi.Message")
	}

	return payloadMessage.Bytes(), nil
}

func (mme *MIDIMessageEncode) Type() string {
	return mme.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "midi.message.encode",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &MIDIMessageEncode{config: config}, nil
		},
	})
}
