package module

import (
	"context"
	"errors"
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
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "mqtt.client",
		New: func(ctx context.Context, config config.ModuleConfig) (Module, error) {
			params := config.Params
			broker, ok := params["broker"]

			if !ok {
				return nil, errors.New("mqtt.client requires a broker parameter")
			}

			brokerString, ok := broker.(string)

			if !ok {
				return nil, errors.New("mqtt.client broker must be string")
			}

			topic, ok := params["topic"]

			if !ok {
				return nil, errors.New("mqtt.client requires a topic parameter")
			}

			topicString, ok := topic.(string)

			if !ok {
				return nil, errors.New("mqtt.client topic must be string")
			}

			clientId, ok := params["clientId"]

			if !ok {
				return nil, errors.New("mqtt.client requires a clientId parameter")
			}

			clientIdString, ok := clientId.(string)

			if !ok {
				return nil, errors.New("mqtt.client clientId must be string")
			}

			router, ok := ctx.Value(route.RouterContextKey).(route.RouteIO)

			if !ok {
				return nil, errors.New("mqtt.client unable to get router from context")
			}

			return &MQTTClient{config: config, Broker: brokerString, Topic: topicString, ClientID: clientIdString, ctx: ctx, router: router, logger: CreateLogger(config)}, nil
		},
	})
}

func (mc *MQTTClient) Id() string {
	return mc.config.Id
}

func (mc *MQTTClient) Type() string {
	return mc.config.Type
}

func (mc *MQTTClient) Run() error {
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
