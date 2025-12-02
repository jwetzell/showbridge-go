package showbridge

import (
	"fmt"
	"log/slog"

	"github.com/jwetzell/showbridge-go/internal/processing"
	"github.com/nats-io/nats.go"
)

type NATSClient struct {
	config  ModuleConfig
	router  *Router
	URL     string
	Subject string
	client  *nats.Conn
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.nats.client",
		New: func(config ModuleConfig) (Module, error) {
			params := config.Params
			url, ok := params["url"]

			if !ok {
				return nil, fmt.Errorf("net.nats.client requires a url parameter")
			}

			urlString, ok := url.(string)

			if !ok {
				return nil, fmt.Errorf("net.nats.client url must be string")
			}

			subject, ok := params["subject"]

			if !ok {
				return nil, fmt.Errorf("net.nats.client requires a subject parameter")
			}

			subjectString, ok := subject.(string)

			if !ok {
				return nil, fmt.Errorf("net.nats.client subject must be string")
			}

			return &NATSClient{config: config, URL: urlString, Subject: subjectString}, nil
		},
	})
}

func (nc *NATSClient) Id() string {
	return nc.config.Id
}

func (nc *NATSClient) Type() string {
	return nc.config.Type
}

func (nc *NATSClient) RegisterRouter(router *Router) {
	nc.router = router
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
			nc.router.HandleInput(nc.config.Id, msg)
		}
	})

	if err != nil {
		return err
	}

	defer sub.Unsubscribe()

	<-nc.router.Context.Done()
	slog.Debug("router context done in module", "id", nc.config.Id)
	return nil
}

func (nc *NATSClient) Output(payload any) error {

	payloadMessage, ok := payload.(processing.NATSMessage)

	if !ok {
		return fmt.Errorf("net.nats.client is only able to output NATSMessage")
	}

	if nc.client == nil {
		return fmt.Errorf("net.nats.client client is not setup")
	}

	if !nc.client.IsConnected() {
		return fmt.Errorf("net.nats.client is not connected")
	}

	err := nc.client.Publish(payloadMessage.Subject, payloadMessage.Payload)

	return err
}
