package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/route"
)

type MQTTClient struct {
	config   config.ModuleConfig
	ctx      context.Context
	router   route.RouteIO
	Broker   string
	ClientID string
	Topic    string
	client   mqtt.Client
	logger   *slog.Logger
	cancel   context.CancelFunc
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "mqtt.client",
		New: func(config config.ModuleConfig) (Module, error) {
			params := config.Params
			brokerString, err := params.GetString("broker")

			if err != nil {
				return nil, fmt.Errorf("mqtt.client broker error: %w", err)
			}

			topicString, err := params.GetString("topic")

			if err != nil {
				return nil, fmt.Errorf("mqtt.client topic error: %w", err)
			}

			clientIdString, err := params.GetString("clientId")

			if err != nil {
				return nil, fmt.Errorf("mqtt.client clientId error: %w", err)
			}

			return &MQTTClient{config: config, Broker: brokerString, Topic: topicString, ClientID: clientIdString, logger: CreateLogger(config)}, nil
		},
	})
}

func (mc *MQTTClient) Id() string {
	return mc.config.Id
}

func (mc *MQTTClient) Type() string {
	return mc.config.Type
}

func (mc *MQTTClient) Start(ctx context.Context) error {
	mc.logger.Debug("running")
	router, ok := ctx.Value(route.RouterContextKey).(route.RouteIO)

	if !ok {
		return errors.New("mqtt.client unable to get router from context")
	}
	mc.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	mc.ctx = moduleContext
	mc.cancel = cancel

	opts := mqtt.NewClientOptions()
	opts.AddBroker(mc.Broker)
	opts.SetClientID(mc.ClientID)
	opts.SetAutoReconnect(true)
	opts.SetCleanSession(false)

	opts.OnConnect = func(c mqtt.Client) {
		token := mc.client.Subscribe(mc.Topic, 1, func(c mqtt.Client, m mqtt.Message) {
			mc.router.HandleInput(mc.ctx, mc.Id(), m)
		})
		token.Wait()
	}

	mc.client = mqtt.NewClient(opts)
	defer mc.client.Disconnect(250)

	token := mc.client.Connect()

	token.Wait()
	err := token.Error()
	if err != nil {
		return err
	}

	<-mc.ctx.Done()
	mc.logger.Debug("done")
	return nil
}

func (mc *MQTTClient) Output(ctx context.Context, payload any) error {
	payloadMessage, ok := payload.(mqtt.Message)

	if !ok {
		return errors.New("mqtt.client is only able to output a MQTTMessage")
	}

	if mc.client == nil {
		return errors.New("mqtt.client client is not setup")
	}

	if !mc.client.IsConnected() {
		return errors.New("mqtt.client is not connected")
	}

	token := mc.client.Publish(payloadMessage.Topic(), payloadMessage.Qos(), payloadMessage.Retained(), payloadMessage.Payload())

	token.Wait()

	return token.Error()
}

func (mc *MQTTClient) Stop() {
	mc.cancel()
}
