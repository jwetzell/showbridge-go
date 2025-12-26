package processor

import (
	"context"
	"fmt"

	"github.com/jwetzell/artnet-go"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type ArtNetPacketEncode struct {
	config config.ProcessorConfig
}

func (ape *ArtNetPacketEncode) Process(ctx context.Context, payload any) (any, error) {
	payloadPacket, ok := payload.(artnet.ArtNetPacket)

	if !ok {
		return nil, fmt.Errorf("artnet.packet.encode processor only accepts an ArtNetPacket")
	}

	payloadBytes, err := payloadPacket.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return payloadBytes, nil
}

func (ape *ArtNetPacketEncode) Type() string {
	return ape.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "artnet.packet.encode",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &ArtNetPacketEncode{config: config}, nil
		},
	})
}
