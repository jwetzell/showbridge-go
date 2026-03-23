package processor

import (
	"context"
	"fmt"

	"github.com/jwetzell/artnet-go"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type ArtNetPacketDecode struct {
	config config.ProcessorConfig
}

func (apd *ArtNetPacketDecode) Process(ctx context.Context, wrappedPayload common.WrappedPayload) (common.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadBytes, ok := common.GetAnyAsByteSlice(payload)

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

func (apd *ArtNetPacketDecode) Type() string {
	return apd.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "artnet.packet.decode",
		Title: "Decode ArtNet Packet",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &ArtNetPacketDecode{config: config}, nil
		},
	})
}
