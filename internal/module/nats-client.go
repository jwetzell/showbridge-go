package module

import (
	"context"
	"fmt"
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
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "nats.client",
		New: func(ctx context.Context, config config.ModuleConfig, router route.RouteIO) (Module, error) {
			params := config.Params
			url, ok := params["url"]

			if !ok {
				return nil, fmt.Errorf("nats.client requires a url parameter")
			}

			urlString, ok := url.(string)

			if !ok {
				return nil, fmt.Errorf("nats.client url must be string")
			}

			subject, ok := params["subject"]

			if !ok {
				return nil, fmt.Errorf("nats.client requires a subject parameter")
			}

			subjectString, ok := subject.(string)

			if !ok {
				return nil, fmt.Errorf("nats.client subject must be string")
			}

			return &NATSClient{config: config, URL: urlString, Subject: subjectString, ctx: ctx, router: router}, nil
		},
	})
}

func (nc *NATSClient) Id() string {
	return nc.config.Id
}

func (nc *NATSClient) Type() string {
	return nc.config.Type
}

func (nc *NATSClient) Run() error {
	client, err := nats.Connect(nc.URL, nats.RetryOnFailedConnect(true))

	if err != nil {
		return err
	}

	nc.client = client

	defer client.Drain()
	defer client.Close()

	sub, err := nc.client.Subscribe(nc.Subject, func(msg *nats.Msg) {
		if nc.router != nil {
			nc.router.HandleInput(nc.Id(), msg)
		}
	})

	if err != nil {
		return err
	}

	defer sub.Unsubscribe()

	<-nc.ctx.Done()
	slog.Debug("router context done in module", "id", nc.Id())
	return nil
}

func (nc *NATSClient) Output(payload any) error {

	payloadMessage, ok := payload.(processor.NATSMessage)

	if !ok {
		return fmt.Errorf("nats.client is only able to output NATSMessage")
	}

	if nc.client == nil {
		return fmt.Errorf("nats.client client is not setup")
	}

	if !nc.client.IsConnected() {
		return fmt.Errorf("nats.client is not connected")
	}

	err := nc.client.Publish(payloadMessage.Subject, payloadMessage.Payload)

	return err
}
