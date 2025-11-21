package showbridge

import (
	"context"
	"fmt"
	"log/slog"
	"net"
)

type UDPClient struct {
	config ModuleConfig
	Host   string
	Port   uint16
	conn   net.Conn
	router *Router
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.udp.client",
		New: func(config ModuleConfig) (Module, error) {
			params := config.Params
			host, ok := params["host"]

			if !ok {
				return nil, fmt.Errorf("net.udp.client requires a host parameter")
			}

			hostString, ok := host.(string)

			if !ok {
				return nil, fmt.Errorf("net.udp.client host must be uint16")
			}

			port, ok := params["port"]
			if !ok {
				return nil, fmt.Errorf("net.udp.client requires a port parameter")
			}

			portNum, ok := port.(float64)

			if !ok {
				return nil, fmt.Errorf("net.udp.client port must be uint16")
			}

			return &UDPClient{Host: hostString, Port: uint16(portNum), config: config}, nil
		},
	})
}

func (uc *UDPClient) Id() string {
	return uc.config.Id
}

func (uc *UDPClient) Type() string {
	return uc.config.Type
}

func (uc *UDPClient) RegisterRouter(router *Router) {
	slog.Debug("registering router", "id", uc.config.Id)
	uc.router = router
}

func (uc *UDPClient) Run(ctx context.Context) error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", uc.Host, uc.Port))
	if err != nil {
		return err
	}
	client, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}

	uc.conn = client
	<-ctx.Done()
	return nil
}

func (uc *UDPClient) Output(payload any) error {
	if uc.conn != nil {
		payloadBytes, ok := payload.([]byte)
		if !ok {
			return fmt.Errorf("net.udp.client is only able to output bytes")
		}
		_, err := uc.conn.Write(payloadBytes)
		return err
	}
	return nil
}
