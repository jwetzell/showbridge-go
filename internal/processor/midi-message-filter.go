//go:build cgo

package processor

import (
	"context"
	"errors"
	"fmt"

	"github.com/jwetzell/showbridge-go/internal/config"
	"gitlab.com/gomidi/midi/v2"
)

type MIDIMessageFilter struct {
	config   config.ProcessorConfig
	MIDIType string
}

func (mmf *MIDIMessageFilter) Process(ctx context.Context, payload any) (any, error) {
	payloadMessage, ok := payload.(midi.Message)

	if !ok {
		return nil, errors.New("midi.message.filter processor only accepts a midi.Message")
	}

	if payloadMessage.Type().String() != mmf.MIDIType {
		return nil, nil
	}

	return payloadMessage, nil
}

func (mmf *MIDIMessageFilter) Type() string {
	return mmf.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "midi.message.filter",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params
			msgTypeString, err := params.GetString("type")
			if err != nil {
				return nil, fmt.Errorf("midi.message.filter type error: %w", err)
			}

			return &MIDIMessageFilter{config: config, MIDIType: msgTypeString}, nil
		},
	})
}
