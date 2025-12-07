package showbridge

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type UDPServer struct {
	Addr   *net.UDPAddr
	config config.ModuleConfig
	router *Router
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.udp.server",
		New: func(config config.ModuleConfig, router *Router) (Module, error) {
			params := config.Params
			port, ok := params["port"]
			if !ok {
				return nil, fmt.Errorf("net.udp.server requires a port parameter")
			}

			portNum, ok := port.(float64)

			if !ok {
				return nil, fmt.Errorf("net.udp.server port must be a number")
			}

			ipString := "0.0.0.0"

			ip, ok := params["ip"]
			if ok {

				specificIpString, ok := ip.(string)

				if !ok {
					return nil, fmt.Errorf("net.udp.server ip must be a string")
				}
				ipString = specificIpString
			}

			addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ipString, uint16(portNum)))
			if err != nil {
				log.Fatalf("error resolving UDP address: %v", err)
			}

			return &UDPServer{Addr: addr, config: config, router: router}, nil
		},
	})
}

func (us *UDPServer) Id() string {
	return us.config.Id
}

func (us *UDPServer) Type() string {
	return us.config.Id
}

func (us *UDPServer) Run() error {

	listener, err := net.ListenUDP("udp", us.Addr)
	if err != nil {
		return err
	}

	defer listener.Close()

	buffer := make([]byte, 1024)
	for {
		select {
		case <-us.router.Context.Done():
			// TODO(jwetzell): cleanup?
			slog.Debug("router context done in module", "id", us.config.Id)
			return nil
		default:
			listener.SetDeadline(time.Now().Add(time.Millisecond * 200))

			numBytes, _, err := listener.ReadFromUDP(buffer)
			if err != nil {
				//NOTE(jwetzell) we hit deadline
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				}
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
