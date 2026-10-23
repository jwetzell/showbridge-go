//go:build cgo

package processor

import (
	"context"
	"errors"

	"github.com/jwetzell/showbridge-go/common"
	"github.com/jwetzell/showbridge-go/config"
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

func (mme *MIDIMessageEncode) Process(ctx context.Context, wrappedPayload common.WrappedPayload, moduleManager common.ModuleManager, inputHandler common.InputHandler) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadMessage, ok := common.GetAnyAs[midi.Message](payload)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("midi.message.encode processor only accepts a midi.Message")
	}

	wrappedPayload.Payload = payloadMessage.Bytes()
	return wrappedPayload, nil
}
