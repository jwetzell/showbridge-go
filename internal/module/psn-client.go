package module

import (
	"context"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/jwetzell/psn-go"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "psn.client",
		Title: "PosiStageNet Client",
		New: func(config config.ModuleConfig) (common.Module, error) {
			return &PSNClient{config: config, decoder: psn.NewDecoder(), logger: CreateLogger(config)}, nil
		},
	})
}

type PSNClient struct {
	config       config.ModuleConfig
	conn         *net.UDPConn
	ctx          context.Context
	inputHandler common.InputHandler
	decoder      *psn.Decoder
	logger       *slog.Logger
	cancel       context.CancelFunc
	connMu       sync.Mutex
}

func (pc *PSNClient) Id() string {
	return pc.config.Id
}

func (pc *PSNClient) Type() string {
	return pc.config.Type
}

func (pc *PSNClient) Start(ctx context.Context, inputHandler common.InputHandler) error {
	pc.logger.Debug("running")
	pc.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	pc.ctx = moduleContext
	pc.cancel = cancel

	addr, err := net.ResolveUDPAddr("udp", "236.10.10.10:56565")
	if err != nil {
		return err
	}

	client, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		return err
	}

	pc.connMu.Lock()
	pc.conn = client
	pc.connMu.Unlock()

	buffer := make([]byte, 2048)
	for {
		select {
		case <-pc.ctx.Done():
			return nil
		default:
			pc.connMu.Lock()
			pc.conn.SetDeadline(time.Now().Add(time.Millisecond * 200))

			numBytes, _, err := pc.conn.ReadFromUDP(buffer)
			pc.connMu.Unlock()
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

				if pc.inputHandler != nil {
					// TODO(jwetzell): better input handling
					for _, tracker := range pc.decoder.Trackers {
						pc.inputHandler(pc.ctx, pc.Id(), tracker)
					}
				} else {
					pc.logger.Error("has no input handler")
				}
			}
		}
	}
}

func (pc *PSNClient) Stop() {
	if pc.cancel != nil {
		pc.cancel()
	}
	pc.connMu.Lock()
	defer pc.connMu.Unlock()
	if pc.conn != nil {
		pc.conn.Close()
		pc.conn = nil
	}
	pc.logger.Debug("done")
}
