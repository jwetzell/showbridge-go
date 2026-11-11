package processor

import (
	"context"
	"errors"

	freeD "github.com/jwetzell/free-d-go"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/internal/types"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:        "freed.decode",
		Title:       "Decode FreeD",
		Description: "Decodes bytes into a FreeD Position",
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			return &FreeDDecode{config: config}, nil
		},
	})
}

type FreeDDecode struct {
	config config.ProcessorConfig
	buf    [29]byte
}

func (fd *FreeDDecode) Process(next routing.ProcessFunc) routing.ProcessFunc {
	return func(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
		payload := wrappedPayload.Payload
		payloadBytes, ok := types.GetAnyAsByteSlice(payload)

		if !ok {
			return wrappedPayload, errors.New("freed.decode processor only accepts a []byte")
		}

		if len(payloadBytes) != 29 {
			return wrappedPayload, errors.New("freeD packet must be exactly 29 bytes")
		}

		copy(fd.buf[:], payloadBytes)

		payloadMessage, err := freeD.Decode(fd.buf)
		if err != nil {
			return wrappedPayload, err
		}
		wrappedPayload.Payload = payloadMessage
		return next(ctx, wrappedPayload, moduleManager, router)
	}
}
