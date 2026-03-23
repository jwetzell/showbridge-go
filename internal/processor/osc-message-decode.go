package processor

import (
	"context"
	"errors"
	"fmt"

	osc "github.com/jwetzell/osc-go"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type OSCMessageDecode struct {
	config config.ProcessorConfig
}

func (omd *OSCMessageDecode) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadBytes, ok := common.GetAnyAsByteSlice(payload)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("osc.message.decode processor only accepts a []byte payload")
	}

	if len(payloadBytes) == 0 {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("osc.message.decode processor can't work on empty []byte")
	}

	if payloadBytes[0] != '/' {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("osc.message.decode processor needs an OSC looking []byte")
	}

	message, err := osc.MessageFromBytes(payloadBytes)
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("osc.message.decode processor failed to decode OSC message: %w", err)
	}
	wrappedPayload.Payload = message
	return wrappedPayload, nil
}

func (omd *OSCMessageDecode) Type() string {
	return omd.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "osc.message.decode",
		Title: "Decode OSC Message",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &OSCMessageDecode{config: config}, nil
		},
	})
}
