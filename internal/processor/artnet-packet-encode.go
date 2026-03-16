package processor

import (
	"context"
	"fmt"

	"github.com/jwetzell/artnet-go"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type ArtNetPacketEncode struct {
	config config.ProcessorConfig
}

func (ape *ArtNetPacketEncode) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadPacket, ok := common.GetAnyAs[artnet.ArtNetPacket](payload)

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
