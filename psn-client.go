package showbridge

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/jwetzell/psn-go"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type PSNClient struct {
	config  config.ModuleConfig
	conn    *net.UDPConn
	router  *Router
	decoder *psn.Decoder
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.psn.client",
		New: func(config config.ModuleConfig) (Module, error) {

			return &PSNClient{config: config, decoder: psn.NewDecoder()}, nil
		},
	})
}

func (pc *PSNClient) Id() string {
	return pc.config.Id
}

func (pc *PSNClient) Type() string {
	return pc.config.Type
}

func (pc *PSNClient) RegisterRouter(router *Router) {
	pc.router = router
}

func (pc *PSNClient) Run() error {

	addr, err := net.ResolveUDPAddr("udp", "236.10.10.10:56565")
	if err != nil {
		return err
	}

	client, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer client.Close()

	pc.conn = client

	buffer := make([]byte, 2048)
	for {
		select {
		case <-pc.router.Context.Done():
			// TODO(jwetzell): cleanup?
			slog.Debug("router context done in module", "id", pc.config.Id)
			return nil
		default:
			pc.conn.SetDeadline(time.Now().Add(time.Millisecond * 200))

			numBytes, _, err := pc.conn.ReadFromUDP(buffer)
			if err != nil {
				//NOTE(jwetzell) we hit deadline
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				}
				return err
			}

			if numBytes > 0 {
				message := buffer[:numBytes]
				err := pc.decoder.Decode(message)
				if err != nil {
					slog.Error("net.psn.client problem decoding psn traffic", "id", pc.config.Id, "error", err)
				}

				if pc.router != nil {
					for _, tracker := range pc.decoder.Trackers {
						pc.router.HandleInput(pc.config.Id, tracker)
					}
				} else {
					slog.Error("net.psn.client has no router", "id", pc.config.Id)
				}
			}
		}
	}
}

func (pc *PSNClient) Output(payload any) error {
	return fmt.Errorf("net.psn.client output is not implemented")
}
