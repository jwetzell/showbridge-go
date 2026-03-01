package processor

import (
	"context"
	"fmt"

	"github.com/jwetzell/artnet-go"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type ArtNetPacketFilter struct {
	config config.ProcessorConfig
	OpCode uint16
}

func (apf *ArtNetPacketFilter) Process(ctx context.Context, payload any) (any, error) {
	payloadPacket, ok := payload.(artnet.ArtNetPacket)

	if !ok {
		return nil, fmt.Errorf("artnet.packet.filter processor only accepts an ArtNetPacket")
	}

	if payloadPacket.GetOpCode() != apf.OpCode {
		return nil, nil
	}

	return payloadPacket, nil
}

func (apf *ArtNetPacketFilter) Type() string {
	return apf.config.Type
}

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type: "artnet.packet.filter",
		New: func(config config.ProcessorConfig) (Processor, error) {
			params := config.Params

			opCodeNum, err := params.GetInt("opCode")
			if err != nil {
				return nil, fmt.Errorf("artnet.packet.filter opCode error: %w", err)
			}

			return &ArtNetPacketFilter{config: config, OpCode: uint16(opCodeNum)}, nil
		},
	})
}
