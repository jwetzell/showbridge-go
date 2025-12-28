package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/jwetzell/psn-go"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/route"
)

type PSNClient struct {
	config  config.ModuleConfig
	conn    *net.UDPConn
	ctx     context.Context
	router  route.RouteIO
	decoder *psn.Decoder
	logger  *slog.Logger
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "psn.client",
		New: func(ctx context.Context, config config.ModuleConfig) (Module, error) {
			router, ok := ctx.Value(route.RouterContextKey).(route.RouteIO)

			if !ok {
				return nil, errors.New("psn.client unable to get router from context")
			}
			return &PSNClient{config: config, decoder: psn.NewDecoder(), ctx: ctx, router: router, logger: CreateLogger(config)}, nil
		},
	})
}

func (pc *PSNClient) Id() string {
	return pc.config.Id
}

func (pc *PSNClient) Type() string {
	return pc.config.Type
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
		case <-pc.ctx.Done():
			// TODO(jwetzell): cleanup?
			pc.logger.Debug("done")
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
					pc.logger.Error("problem decoding psn traffic", "error", err)
				}

				if pc.router != nil {
					for _, tracker := range pc.decoder.Trackers {
						pc.router.HandleInput(pc.Id(), tracker)
					}
				} else {
					pc.logger.Error("has no router")
				}
			}
		}
	}
}

func (pc *PSNClient) Output(ctx context.Context, payload any) error {
	return fmt.Errorf("psn.client output is not implemented")
}
