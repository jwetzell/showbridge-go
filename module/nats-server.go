package module

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/params"
	"github.com/jwetzell/showbridge-go/routing"
	"github.com/nats-io/nats-server/v2/server"
)

func init() {
	RegisterModule(ModuleRegistration{
		Type:        "nats.server",
		Title:       "NATS Server",
		Description: "Starts a NATS server on the specified IP and port",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"ip": {
					Title:       "IP",
					Description: "the IP address to bind the NATS server to",
					Type:        "string",
					Default:     json.RawMessage(`"0.0.0.0"`),
				},
				"port": {
					Title:       "Port",
					Description: "the port for the NATS server to listen on",
					Type:        "integer",
					Minimum:     jsonschema.Ptr[float64](1024),
					Maximum:     jsonschema.Ptr[float64](65535),
					Default:     json.RawMessage(`4222`),
				},
			},
			Required:             []string{},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(moduleConfig config.ModuleConfig) (routing.Module, error) {
			moduleParams := moduleConfig.Params
			portNum, err := moduleParams.GetInt("port")
			if err != nil {
				if errors.Is(err, params.ErrParamNotFound) {
					portNum = 4222
				} else {
					return nil, fmt.Errorf("nats.server port error: %w", err)
				}
			}

			ipString, err := moduleParams.GetString("ip")
			if err != nil {
				if errors.Is(err, params.ErrParamNotFound) {
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

type NATSServer struct {
	config   config.ModuleConfig
	ctx      context.Context
	Ip       string
	Port     int
	router   routing.Router
	server   *server.Server
	logger   *slog.Logger
	cancel   context.CancelFunc
	serverMu sync.Mutex
}

func (ns *NATSServer) Start(ctx context.Context, router routing.Router) error {
	ns.logger.Debug("running")
	ns.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	ns.ctx = moduleContext
	ns.cancel = cancel

	natsServer, err := server.NewServer(&server.Options{
		Host:   ns.Ip,
		Port:   ns.Port,
		NoLog:  true,
		NoSigs: true,
	})

	if err != nil {
		return err
	}

	ns.serverMu.Lock()
	ns.server = natsServer
	natsServer.Start()

	if !natsServer.ReadyForConnections(5 * time.Second) {
		ns.serverMu.Unlock()
		return errors.New("nats.server failed to start")
	}
	ns.serverMu.Unlock()

	<-ns.ctx.Done()
	ns.logger.Debug("done")
	return nil
}

func (ns *NATSServer) Stop() {
	if ns.cancel != nil {
		defer ns.cancel()
	}
	ns.serverMu.Lock()
	defer ns.serverMu.Unlock()
	if ns.server != nil {
		ns.server.Shutdown()
		ns.server.WaitForShutdown()
	}
}
