package showbridge

import (
	"context"
	"fmt"
)

type Router struct {
	Context           context.Context
	ProtocolInstances []Protocol
}

func NewRouter(ctx context.Context, config Config) (*Router, error) {

	router := Router{
		Context:           ctx,
		ProtocolInstances: []Protocol{},
	}

	for _, protocolDecl := range config.Protocols {

		protocolInfo, ok := protocolRegistry[protocolDecl.Type]
		if !ok {
			return nil, fmt.Errorf("problem loading protocol registration for protocol type: %s", protocolDecl.Type)
		}

		protocolInstance, err := protocolInfo.New(protocolDecl.Params)
		if err != nil {
			return nil, err
		}

		router.ProtocolInstances = append(router.ProtocolInstances, protocolInstance)

	}

	return &router, nil
}

func (r *Router) Run() {
	for _, protocolInstance := range r.ProtocolInstances {
		protocolInstance.Run(r.Context)
	}
}
