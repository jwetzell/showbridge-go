package module

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type MQTTClient struct {
	config   config.ModuleConfig
	ctx      context.Context
	router   common.RouteIO
	Broker   string
	ClientID string
	Topic    string
	QoS      byte
	Retained bool
	client   mqtt.Client
	logger   *slog.Logger
	cancel   context.CancelFunc
	clientMu sync.Mutex
}

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "mqtt.client",
		Title: "MQTT Client",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"broker": {
					Title: "Broker URL",
					Type:  "string",
				},
				"topic": {
					Title: "Topic",
					Type:  "string",
				},
				"clientId": {
					Title: "Client ID",
					Type:  "string",
				},
				"qos": {
					Title:   "QoS",
					Type:    "integer",
					Minimum: jsonschema.Ptr[float64](0),
					Maximum: jsonschema.Ptr[float64](2),
					Default: json.RawMessage(`0`),
				},
				"retained": {
					Title:   "Retained",
					Type:    "boolean",
					Default: json.RawMessage(`false`),
				},
			},
			Required:             []string{"broker", "topic", "clientId"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(moduleConfig config.ModuleConfig) (common.Module, error) {
			params := moduleConfig.Params
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

			qosString, err := params.GetInt("qos")

			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					qosString = 0
				} else {
					return nil, fmt.Errorf("mqtt.client qos error: %w", err)
				}
			}

			retainedBool, err := params.GetBool("retained")

			if err != nil {
				if errors.Is(err, config.ErrParamNotFound) {
					retainedBool = false
				} else {
					return nil, fmt.Errorf("mqtt.client retained error: %w", err)
				}
			}

			return &MQTTClient{config: moduleConfig, Broker: brokerString, Topic: topicString, ClientID: clientIdString, QoS: byte(qosString), Retained: retainedBool, logger: CreateLogger(moduleConfig)}, nil
		},
	})
}

func (mc *MQTTClient) Id() string {
	return mc.config.Id
}

func (mc *MQTTClient) Type() string {
	return mc.config.Type
}

func (mc *MQTTClient) Start(ctx context.Context, router common.RouteIO) error {
	mc.logger.Debug("running")
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

	mc.clientMu.Lock()
	mc.client = mqtt.NewClient(opts)

	token := mc.client.Connect()

	token.Wait()
	err := token.Error()
	if err != nil {
		return err
	}
	mc.clientMu.Unlock()

	<-mc.ctx.Done()
	return nil
}

func (mc *MQTTClient) Publish(ctx context.Context, topic string, payload any) error {
	payloadBytes, ok := common.GetAnyAsByteSlice(payload)

	if !ok {
		payloadString, ok := common.GetAnyAs[string](payload)
		if !ok {
			return errors.New("mqtt.client is only able to publish bytes or string")
		}
		payloadBytes = []byte(payloadString)
	}

	if mc.client == nil {
		return errors.New("mqtt.client client is not setup")
	}

	if !mc.client.IsConnected() {
		return errors.New("mqtt.client is not connected")
	}

	token := mc.client.Publish(topic, mc.QoS, mc.Retained, payloadBytes)

	token.Wait()

	return token.Error()
}

func (mc *MQTTClient) Stop() {
	if mc.cancel != nil {
		mc.cancel()
	}
	mc.clientMu.Lock()
	defer mc.clientMu.Unlock()
	if mc.client != nil {
		mc.client.Disconnect(250)
		mc.client = nil
	}
	mc.logger.Debug("done")
}
