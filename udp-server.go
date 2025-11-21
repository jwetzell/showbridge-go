package showbridge

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
)

type UDPServer struct {
	Port   uint16
	config ModuleConfig
	router *Router
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.udp.server",
		New: func(config ModuleConfig) (Module, error) {
			params := config.Params
			port, ok := params["port"]
			if !ok {
				return nil, fmt.Errorf("net.udp.server requires a port parameter")
			}

			portNum, ok := port.(float64)

			if !ok {
				return nil, fmt.Errorf("net.udp.server port must be uint16")
			}

			return &UDPServer{Port: uint16(portNum), config: config}, nil
		},
	})
}

func (us *UDPServer) Id() string {
	return us.config.Id
}

func (us *UDPServer) Type() string {
	return us.config.Id
}

func (us *UDPServer) RegisterRouter(router *Router) {
	us.router = router
}

func (us *UDPServer) Run(ctx context.Context) error {

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", us.Port))
	if err != nil {
		log.Fatalf("error resolving UDP address: %v", err)
	}

	listener, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	defer listener.Close()

	buffer := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			numBytes, _, err := listener.ReadFromUDP(buffer)
			if err != nil {
				return err
			}
			message := buffer[:numBytes]
			if us.router != nil {
				us.router.HandleInput(us.config.Id, message)
			} else {
				slog.Error("net.udp.server has no router", "id", us.config.Id)
			}
		}
	}

}

func (us *UDPServer) Output(payload any) error {
	return fmt.Errorf("net.udp.server output is not implemented")
}
