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
		Type:  "midi.message.decode",
		Title: "Decode MIDI Message",
		New: func(config config.ProcessorConfig) (common.Processor, error) {
			return &MIDIMessageDecode{config: config}, nil
		},
	})
}

type MIDIMessageDecode struct {
	config config.ProcessorConfig
}

func (mmd *MIDIMessageDecode) Process(ctx context.Context, wrappedPayload common.WrappedPayload, moduleManager common.ModuleManager, inputHandler common.InputHandler) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadBytes, ok := common.GetAnyAsByteSlice(payload)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("midi.message.decode processor only accepts a []byte")
	}

	payloadMessage := midi.Message(payloadBytes)

	wrappedPayload.Payload = payloadMessage
	return wrappedPayload, nil
}
