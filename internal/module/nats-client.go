package module

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/nats-io/nats.go"
)

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "nats.client",
		Title: "NATS Client",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"url": {
					Title:       "NATS Server URL",
					Description: "the URL of the NATS server to connect to",
					Type:        "string",
				},
				"subject": {
					Title:       "Subject",
					Description: "the subject to subscribe to",
					Type:        "string",
				},
			},
			Required:             []string{"url", "subject"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ModuleConfig) (common.Module, error) {
			params := config.Params
			urlString, err := params.GetString("url")
			if err != nil {
				return nil, errors.New("nats.client url error: " + err.Error())
			}

			subjectString, err := params.GetString("subject")

			if err != nil {
				return nil, errors.New("nats.client subject error: " + err.Error())
			}

			return &NATSClient{config: config, URL: urlString, Subject: subjectString, logger: CreateLogger(config)}, nil
		},
	})
}

type NATSClient struct {
	config       config.ModuleConfig
	ctx          context.Context
	inputHandler common.InputHandler
	URL          string
	Subject      string
	client       *nats.Conn
	logger       *slog.Logger
	cancel       context.CancelFunc
	sub          *nats.Subscription
	subMu        sync.Mutex
	clientMu     sync.Mutex
}

func (nc *NATSClient) Id() string {
	return nc.config.Id
}

func (nc *NATSClient) Type() string {
	return nc.config.Type
}

func (nc *NATSClient) Start(ctx context.Context, inputHandler common.InputHandler) error {
	nc.logger.Debug("running")
	nc.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	nc.ctx = moduleContext
	nc.cancel = cancel

	client, err := nats.Connect(nc.URL, nats.RetryOnFailedConnect(true))

	if err != nil {
		return err
	}

	nc.clientMu.Lock()
	nc.client = client
	nc.clientMu.Unlock()

	sub, err := nc.client.Subscribe(nc.Subject, func(msg *nats.Msg) {
		if nc.inputHandler != nil {
			nc.inputHandler(nc.ctx, nc.Id(), msg)
		}
	})

	if err != nil {
		return err
	}
	nc.subMu.Lock()
	nc.sub = sub
	nc.subMu.Unlock()

	<-nc.ctx.Done()
	nc.logger.Debug("done")
	return nil
}

func (nc *NATSClient) Publish(ctx context.Context, topic string, payload any) error {

	payloadBytes, ok := common.GetAnyAsByteSlice(payload)

	if !ok {
		payloadString, ok := common.GetAnyAs[string](payload)
		if !ok {
			return errors.New("nats.client is only able to publish bytes or string")
		}
		payloadBytes = []byte(payloadString)
	}

	nc.clientMu.Lock()
	defer nc.clientMu.Unlock()

	if nc.client == nil {
		return errors.New("nats.client client is not setup")
	}

	if !nc.client.IsConnected() {
		return errors.New("nats.client is not connected")
	}

	err := nc.client.Publish(topic, payloadBytes)

	return err
}

func (nc *NATSClient) Stop() {
	if nc.cancel != nil {
		defer nc.cancel()
	}
	nc.subMu.Lock()
	defer nc.subMu.Unlock()
	if nc.sub != nil {
		nc.sub.Unsubscribe()
	}

	nc.clientMu.Lock()
	defer nc.clientMu.Unlock()
	if nc.client != nil {
		nc.client.Drain()
		// TODO(jwetzell): setup closed callback to get when client is fully closed
		nc.client.Close()
	}
}
