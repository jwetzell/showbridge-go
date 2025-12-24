package processor

import (
	"context"
	"fmt"

	"github.com/jwetzell/artnet-go"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type ArtNetPacketDecode struct {
	config config.ProcessorConfig
}

func (apd *ArtNetPacketDecode) Process(ctx context.Context, payload any) (any, error) {
	payloadBytes, ok := payload.([]byte)

	if !ok {
		return nil, fmt.Errorf("artnet.packet.decode processor only accepts a []byte")
	}

	payloadMessage, err := artnet.Decode(payloadBytes)

	if err != nil {
		return nil, err
	}

	return payloadMessage, nil
}

func (apd *ArtNetPacketDecode) Type() string {
	return apd.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "artnet.packet.decode",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &ArtNetPacketDecode{config: config}, nil
		},
	})
}
