package processor

import (
	"context"
	"fmt"

	"github.com/jwetzell/artnet-go"
	"github.com/jwetzell/showbridge-go/common"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/payload"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "artnet.packet.encode",
		Title: "Encode ArtNet Packet",
		New: func(config config.ProcessorConfig) (common.Processor, error) {
			return &ArtNetPacketEncode{config: config}, nil
		},
	})
}

type ArtNetPacketEncode struct {
	config config.ProcessorConfig
}

func (ape *ArtNetPacketEncode) Process(ctx context.Context, wrappedPayload payload.WrappedPayload, moduleManager common.ModuleManager, inputHandler common.InputHandler) (payload.WrappedPayload, error) {
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
