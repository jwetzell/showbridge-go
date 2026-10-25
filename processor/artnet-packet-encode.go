package processor

import (
	"context"
	"fmt"

	"github.com/jwetzell/artnet-go"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "artnet.packet.encode",
		Title: "Encode ArtNet Packet",
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			return &ArtNetPacketEncode{config: config}, nil
		},
	})
}

type ArtNetPacketEncode struct {
	config config.ProcessorConfig
}

func (ape *ArtNetPacketEncode) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadPacket, ok := payload.(artnet.ArtNetPacket)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, fmt.Errorf("artnet.packet.encode processor only accepts an ArtNetPacket")
	}

	payloadBytes, err := payloadPacket.MarshalBinary()
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}

	wrappedPayload.Payload = payloadBytes

	return wrappedPayload, nil
}
