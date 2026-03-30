//go:build cgo || js

package processor

import (
	"context"
	"errors"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"gitlab.com/gomidi/midi/v2"
)

type MIDIMessageDecode struct {
	config config.ProcessorConfig
}

func (mmd *MIDIMessageDecode) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
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

func (mmd *MIDIMessageDecode) Type() string {
	return mmd.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "midi.message.decode",
		Title: "Decode MIDI Message",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &MIDIMessageDecode{config: config}, nil
		},
	})
}
