//go:build cgo

package processor

import (
	"context"
	"errors"

	"github.com/jwetzell/showbridge-go/common"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/payload"
	"gitlab.com/gomidi/midi/v2"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "midi.message.encode",
		Title: "Encode MIDI Message",
		New: func(config config.ProcessorConfig) (common.Processor, error) {
			return &MIDIMessageEncode{config: config}, nil
		},
	})
}

type MIDIMessageEncode struct {
	config config.ProcessorConfig
}

func (mme *MIDIMessageEncode) Process(ctx context.Context, wrappedPayload payload.WrappedPayload, moduleManager common.ModuleManager, inputHandler common.InputHandler) (payload.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadMessage, ok := payload.(midi.Message)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("midi.message.encode processor only accepts a midi.Message")
	}

	wrappedPayload.Payload = payloadMessage.Bytes()
	return wrappedPayload, nil
}
