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

			opCode, ok := params["opCode"]
			if !ok {
				return nil, fmt.Errorf("artnet.packet.filter requires an opCode parameter")
			}
			opCodeNum, ok := opCode.(float64)
			if !ok {
				return nil, fmt.Errorf("artnet.packet.filter opCode must be a number")
			}

			return &ArtNetPacketFilter{config: config, OpCode: uint16(opCodeNum)}, nil
		},
	})
}
