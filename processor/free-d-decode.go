package processor

import (
	"context"
	"errors"

	freeD "github.com/jwetzell/free-d-go"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/params"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterProcessor(ProcessorRegistration{
		Type:  "freed.decode",
		Title: "Decode FreeD",
		New: func(config config.ProcessorConfig) (routing.Processor, error) {
			return &FreeDDecode{config: config}, nil
		},
	})
}

type FreeDDecode struct {
	config config.ProcessorConfig
	buf    [29]byte
}

func (fd *FreeDDecode) Process(ctx context.Context, wrappedPayload routing.WrappedPayload, moduleManager routing.ModuleManager, router routing.Router) (routing.WrappedPayload, error) {
	payload := wrappedPayload.Payload
	payloadBytes, ok := params.GetAnyAsByteSlice(payload)

	if !ok {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("freed.decode processor only accepts a []byte")
	}

	if len(payloadBytes) != 29 {
		wrappedPayload.End = true
		return wrappedPayload, errors.New("freeD packet must be exactly 29 bytes")
	}

	copy(fd.buf[:], payloadBytes)

	payloadMessage, err := freeD.Decode(fd.buf)
	if err != nil {
		wrappedPayload.End = true
		return wrappedPayload, err
	}
	wrappedPayload.Payload = payloadMessage
	return wrappedPayload, nil
}
