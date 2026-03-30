//go:build cgo || js

package processor

import (
	"context"
	"errors"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"gitlab.com/gomidi/midi/v2"
)

type MIDIMessageEncode struct {
	config config.ProcessorConfig
}

func (mme *MIDIMessageEncode) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadMessage, ok := common.GetAnyAs[midi.Message](payload)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("midi.message.encode processor only accepts a midi.Message")
	}

	wrappedPayload.Payload = payloadMessage.Bytes()
	return wrappedPayload, nil
}

func (mme *MIDIMessageEncode) Type() string {
	return mme.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "midi.message.encode",
		Title: "Encode MIDI Message",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &MIDIMessageEncode{config: config}, nil
		},
	})
}
