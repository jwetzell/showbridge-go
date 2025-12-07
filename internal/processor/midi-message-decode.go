//go:build cgo

package processor

import (
	"context"
	"fmt"

	"github.com/jwetzell/showbridge-go/internal/config"
	"gitlab.com/gomidi/midi/v2"
)

type MIDIMessageDecode struct {
	config config.ProcessorConfig
}

func (mmd *MIDIMessageDecode) Process(ctx context.Context, payload any) (any, error) {
	payloadBytes, ok := payload.([]byte)

	if !ok {
		return nil, fmt.Errorf("midi.message.decode processor only accepts a []byte")
	}

	payloadMessage := midi.Message(payloadBytes)

	return payloadMessage, nil
}

func (mmd *MIDIMessageDecode) Type() string {
	return mmd.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "midi.message.decode",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &MIDIMessageDecode{config: config}, nil
		},
	})
}
