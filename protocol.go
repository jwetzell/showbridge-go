package showbridge

import (
	"context"
	"fmt"
	"sync"
)

type Protocol interface {
	Run(context.Context) error
}

type ProtocolConfig struct {
	Type   string         `json:"type"`
	Params map[string]any `json:"params"`
}

type ProtocolRegistration struct {
	Type string `json:"type"`
	New  func(map[string]any) (Protocol, error)
}

func RegisterProtocol(proto ProtocolRegistration) {

	if proto.Type == "" {
		panic("protocol type is missing")
	}
	if proto.New == nil {
		panic("missing ProtocolInfo.New")
	}

	protocolRegistryMu.Lock()
	defer protocolRegistryMu.Unlock()

	if _, ok := protocolRegistry[string(proto.Type)]; ok {
		panic(fmt.Sprintf("protocol already registered: %s", proto.Type))
	}
	protocolRegistry[string(proto.Type)] = proto
}

var (
	protocolRegistryMu sync.RWMutex
	protocolRegistry   = make(map[string]ProtocolRegistration)
)
