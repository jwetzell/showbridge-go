package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/route"
	"github.com/nats-io/nats-server/v2/server"
)

type NATSServer struct {
	config config.ModuleConfig
	ctx    context.Context
	Ip     string
	Port   int
	router route.RouteIO
	server *server.Server
	logger *slog.Logger
	cancel context.CancelFunc
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "nats.server",
		New: func(moduleConfig config.ModuleConfig) (Module, error) {
			params := moduleConfig.Params
			portNum, err := params.GetInt("port")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					portNum = 4222
				} else {
					return nil, fmt.Errorf("nats.server port error: %w", err)
				}
			}

			ipString, err := params.GetString("ip")
			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					ipString = "0.0.0.0"
				} else {
					return nil, fmt.Errorf("nats.server ip error: %w", err)
				}
			}

			_, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ipString, uint16(portNum)))
			if err != nil {
				return nil, err
			}
			return &NATSServer{config: moduleConfig, logger: CreateLogger(moduleConfig), Ip: ipString, Port: portNum}, nil
		},
	})
}

func (ns *NATSServer) Id() string {
	return ns.config.Id
}

func (ns *NATSServer) Type() string {
	return ns.config.Type
}

func (ns *NATSServer) Start(ctx context.Context) error {
	ns.logger.Debug("running")
	router, ok := ctx.Value(route.RouterContextKey).(route.RouteIO)

	if !ok {
		return errors.New("nats.server unable to get router from context")
	}

	ns.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	ns.ctx = moduleContext
	ns.cancel = cancel

	natsServer, err := server.NewServer(&server.Options{
		Host:  ns.Ip,
		Port:  ns.Port,
		NoLog: true,
	})

	if err != nil {
		return err
	}

	ns.server = natsServer
	natsServer.Start()
	defer natsServer.Shutdown()

	if !natsServer.ReadyForConnections(5 * time.Second) {
		return errors.New("nats.server failed to start")
	}
	ns.logger.Info("NATS server started", "client_url", natsServer.ClientURL())

	<-ns.ctx.Done()

	ns.logger.Debug("done")
	return nil
}

func (ns *NATSServer) Output(ctx context.Context, payload any) error {
	return errors.ErrUnsupported
}

func (ns *NATSServer) Stop() {
	ns.cancel()
	if ns.server != nil {
		ns.server.Shutdown()
	}
}
