//go:build cgo

package processor

import (
	"context"
	"errors"

	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/params"
	"github.com/jwetzell/showbridge-go/routing"
	"gitlab.com/gomidi/midi/v2"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:        "midi.message.decode",
		Title:       "Decode MIDI Message",
		Description: "Decodes bytes into a MIDI message",
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			return &MIDIMessageDecode{config: config}, nil
		},
	})
}

type MIDIMessageDecode struct {
	config config.ProcessorConfig
}

func (mmd *MIDIMessageDecode) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadBytes, ok := params.GetAnyAsByteSlice(payload)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("midi.message.decode processor only accepts a []byte")
	}

	payloadMessage := midi.Message(payloadBytes)

	wrappedPayload.Payload = payloadMessage
	return wrappedPayload, nil
}
