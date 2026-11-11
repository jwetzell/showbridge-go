package processor

import (
	"context"
	"fmt"

	"github.com/jwetzell/artnet-go"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/types"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:        "artnet.packet.decode",
		Title:       "Decode ArtNet Packet",
		Description: "Decodes an ArtNet packet from a byte slice",
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			return &ArtNetPacketDecode{config: config}, nil
		},
	})
}

type ArtNetPacketDecode struct {
	config config.ProcessorConfig
}

func (apd *ArtNetPacketDecode) Process(next routing.ProcessFunc) routing.ProcessFunc {
	return func(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
		payload := wrappedPayload.Payload
		payloadBytes, ok := types.GetAnyAsByteSlice(payload)

		if !ok {
			return wrappedPayload, fmt.Errorf("artnet.packet.decode processor only accepts a []byte")
		}

		payloadMessage, err := artnet.Decode(payloadBytes)

		if err != nil {
			return wrappedPayload, err
		}

		wrappedPayload.Payload = payloadMessage

		return next(ctx, wrappedPayload, moduleManager, router)
	}
}
