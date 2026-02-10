package module

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processor"
	"github.com/jwetzell/showbridge-go/internal/route"
	"github.com/nats-io/nats.go"
)

type NATSClient struct {
	config  config.ModuleConfig
	ctx     context.Context
	router  route.RouteIO
	URL     string
	Subject string
	client  *nats.Conn
	logger  *slog.Logger
	cancel  context.CancelFunc
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "nats.client",
		New: func(config config.ModuleConfig) (Module, error) {
			params := config.Params
			url, ok := params["url"]

			if !ok {
				return nil, errors.New("nats.client requires a url parameter")
			}

			urlString, ok := url.(string)

			if !ok {
				return nil, errors.New("nats.client url must be a string")
			}

			subject, ok := params["subject"]

			if !ok {
				return nil, errors.New("nats.client requires a subject parameter")
			}

			subjectString, ok := subject.(string)

			if !ok {
				return nil, errors.New("nats.client subject must be a string")
			}

			return &NATSClient{config: config, URL: urlString, Subject: subjectString, logger: CreateLogger(config)}, nil
		},
	})
}

func (nc *NATSClient) Id() string {
	return nc.config.Id
}

func (nc *NATSClient) Type() string {
	return nc.config.Type
}

func (nc *NATSClient) Start(ctx context.Context) error {
	nc.logger.Debug("running")
	router, ok := ctx.Value(route.RouterContextKey).(route.RouteIO)

	if !ok {
		return errors.New("nats.client unable to get router from context")
	}

	nc.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	nc.ctx = moduleContext
	nc.cancel = cancel

	client, err := nats.Connect(nc.URL, nats.RetryOnFailedConnect(true))

	if err != nil {
		return err
	}

	nc.client = client

	defer client.Drain()
	defer client.Close()

	sub, err := nc.client.Subscribe(nc.Subject, func(msg *nats.Msg) {
		if nc.router != nil {
			nc.router.HandleInput(nc.ctx, nc.Id(), msg)
		}
	})

	if err != nil {
		return err
	}

	defer sub.Unsubscribe()

	<-nc.ctx.Done()
	nc.logger.Debug("done")
	return nil
}

func (nc *NATSClient) Output(ctx context.Context, payload any) error {

	payloadMessage, ok := payload.(processor.NATSMessage)

	if !ok {
		return errors.New("nats.client is only able to output NATSMessage")
	}

	if nc.client == nil {
		return errors.New("nats.client client is not setup")
	}

	if !nc.client.IsConnected() {
		return errors.New("nats.client is not connected")
	}

	err := nc.client.Publish(payloadMessage.Subject, payloadMessage.Payload)

	return err
}

func (nc *NATSClient) Stop() {
	nc.cancel()
}
