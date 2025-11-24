package showbridge

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/jwetzell/psn-go"
)

type PSNClient struct {
	config  ModuleConfig
	conn    *net.UDPConn
	router  *Router
	decoder *psn.Decoder
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.psn.client",
		New: func(config ModuleConfig) (Module, error) {

			return &PSNClient{config: config, decoder: psn.NewDecoder()}, nil
		},
	})
}

func (uc *PSNClient) Id() string {
	return uc.config.Id
}

func (uc *PSNClient) Type() string {
	return uc.config.Type
}

func (uc *PSNClient) RegisterRouter(router *Router) {
	uc.router = router
}

func (uc *PSNClient) Run() error {

	addr, err := net.ResolveUDPAddr("udp", "236.10.10.10:56565")
	if err != nil {
		return err
	}

	client, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer client.Close()

	uc.conn = client

	buffer := make([]byte, 2048)
	for {
		select {
		case <-uc.router.Context.Done():
			// TODO(jwetzell): cleanup?
			slog.Debug("router context done in module", "id", uc.config.Id)
			return nil
		default:
			uc.conn.SetDeadline(time.Now().Add(time.Millisecond * 200))

			numBytes, _, err := uc.conn.ReadFromUDP(buffer)
			if err != nil {
				//NOTE(jwetzell) we hit deadline
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				}
				return err
			}

			if numBytes > 0 {
				message := buffer[:numBytes]
				err := uc.decoder.Decode(message)
				if err != nil {
					slog.Error("net.psn.client problem decoding psn traffic", "id", uc.config.Id, "error", err)
				}

				if uc.router != nil {
					for _, tracker := range uc.decoder.Trackers {
						uc.router.HandleInput(uc.config.Id, tracker)
					}
				} else {
					slog.Error("net.psn.client has no router", "id", uc.config.Id)
				}
			}
		}
	}
}

func (uc *PSNClient) Output(payload any) error {
	return fmt.Errorf("net.psn.client output is not implemented")
}
