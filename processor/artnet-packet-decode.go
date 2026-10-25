package processor

import (
	"context"
	"fmt"

	"github.com/jwetzell/artnet-go"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/params"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "artnet.packet.decode",
		Title: "Decode ArtNet Packet",
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			return &ArtNetPacketDecode{config: config}, nil
		},
	})
}

type ArtNetPacketDecode struct {
	config config.ProcessorConfig
}

func (apd *ArtNetPacketDecode) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadBytes, ok := params.GetAnyAsByteSlice(payload)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("artnet.packet.decode processor only accepts a []byte")
	}

	payloadMessage, err := artnet.Decode(payloadBytes)

	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	wrappedPayload.Payload = payloadMessage

	return wrappedPayload, nil
}
