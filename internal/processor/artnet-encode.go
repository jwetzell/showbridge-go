package processor

import (
	"context"
	"fmt"

	"github.com/jwetzell/artnet-go"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type ArtNetEncode struct {
	config config.ProcessorConfig
}

func (ad *ArtNetEncode) Process(ctx context.Context, payload any) (any, error) {
	payloadPacket, ok := payload.(artnet.ArtNetPacket)

	if !ok {
		return nil, fmt.Errorf("artnet.encode processor only accepts an ArtNetPacket")
	}

	payloadBytes, err := payloadPacket.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return payloadBytes, nil
}

func (ad *ArtNetEncode) Type() string {
	return ad.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "artnet.encode",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &ArtNetEncode{config: config}, nil
		},
	})
}
