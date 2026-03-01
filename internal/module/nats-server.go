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
		New: func(config config.ModuleConfig) (Module, error) {
			params := config.Params
			portNum := 4222

			port, ok := params["port"]
			if ok {
				specificportNum, ok := port.(int)
				if !ok {
					specificportNum, ok := port.(float64)
					if !ok {
						return nil, errors.New("nats.server port must be a number")
					}
					portNum = int(specificportNum)
				} else {
					portNum = int(specificportNum)
				}
			}

			ipString := "0.0.0.0"

			ip, ok := params["ip"]
			if ok {

				specificIpString, ok := ip.(string)

				if !ok {
					return nil, errors.New("nats.server ip must be a string")
				}
				ipString = specificIpString
			}

			_, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ipString, uint16(portNum)))
			if err != nil {
				return nil, err
			}
			return &NATSServer{config: config, logger: CreateLogger(config), Ip: ipString, Port: portNum}, nil
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
