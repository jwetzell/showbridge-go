package showbridge

import (
	"context"
	"fmt"
	"log"
	"net"
)

type UDPServer struct {
	Port uint16
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.udp.server",
		New: func(params map[string]any) (Module, error) {
			port, ok := params["port"]
			if !ok {
				return nil, fmt.Errorf("udp server requires a port parameter")
			}

			portNum, ok := port.(float64)

			if !ok {
				return nil, fmt.Errorf("udp server port must be uint16")
			}

			return UDPServer{Port: uint16(portNum)}, nil
		},
	})
}

func (us UDPServer) Run(ctx context.Context) error {

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
			fmt.Println(buffer[:numBytes])
		}
	}

}
