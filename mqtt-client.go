package showbridge

import (
	"fmt"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/processing"
)

type MQTTClient struct {
	config   config.ModuleConfig
	router   *Router
	Broker   string
	ClientID string
	Topic    string
	client   mqtt.Client
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.mqtt.client",
		New: func(config config.ModuleConfig, router *Router) (Module, error) {
			params := config.Params
			broker, ok := params["broker"]

			if !ok {
				return nil, fmt.Errorf("net.mqtt.client requires a broker parameter")
			}

			brokerString, ok := broker.(string)

			if !ok {
				return nil, fmt.Errorf("net.mqtt.client broker must be string")
			}

			topic, ok := params["topic"]

			if !ok {
				return nil, fmt.Errorf("net.mqtt.client requires a topic parameter")
			}

			topicString, ok := topic.(string)

			if !ok {
				return nil, fmt.Errorf("net.mqtt.client topic must be string")
			}

			clientId, ok := params["clientId"]

			if !ok {
				return nil, fmt.Errorf("net.mqtt.client requires a clientId parameter")
			}

			clientIdString, ok := clientId.(string)

			if !ok {
				return nil, fmt.Errorf("net.mqtt.client clientId must be string")
			}

			return &MQTTClient{config: config, Broker: brokerString, Topic: topicString, ClientID: clientIdString, router: router}, nil
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
			mc.router.HandleInput(mc.config.Id, m)
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

	<-mc.router.Context.Done()
	slog.Debug("router context done in module", "id", mc.config.Id)
	return nil
}

func (mc *MQTTClient) Output(payload any) error {
	payloadMessage, ok := payload.(processing.MQTTMessage)

	if !ok {
		return fmt.Errorf("net.mqtt.client is only able to output a MQTTMessage")
	}

	if mc.client == nil {
		return fmt.Errorf("net.mqtt.client client is not setup")
	}

	if !mc.client.IsConnected() {
		return fmt.Errorf("net.mqtt.client is not connected")
	}

	token := mc.client.Publish(payloadMessage.Topic, payloadMessage.QoS, payloadMessage.Retained, payloadMessage.Payload)

	token.Wait()

	return token.Error()
}
