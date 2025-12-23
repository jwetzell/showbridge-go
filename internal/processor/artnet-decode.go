package processor

import (
	"context"
	"fmt"

	"github.com/jwetzell/artnet-go"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type ArtNetDecode struct {
	config config.ProcessorConfig
}

func (ad *ArtNetDecode) Process(ctx context.Context, payload any) (any, error) {
	payloadBytes, ok := payload.([]byte)

	if !ok {
		return nil, fmt.Errorf("artnet.decode processor only accepts a []byte")
	}

	payloadMessage, err := artnet.Decode(payloadBytes)

	if err != nil {
		return nil, err
	}

	return payloadMessage, nil
}

func (ad *ArtNetDecode) Type() string {
	return ad.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "artnet.decode",
		New: func(config config.ProcessorConfig) (Processor, error) {
			return &ArtNetDecode{config: config}, nil
		},
	})
}
