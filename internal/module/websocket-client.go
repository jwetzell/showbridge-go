package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/gorilla/websocket"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

type WebSocketClient struct {
	config config.ModuleConfig
	URL    url.URL
	ctx    context.Context
	conn   *websocket.Conn
	router common.RouteIO
	logger *slog.Logger
	cancel context.CancelFunc
}

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "websocket.client",
		Title: "WebSocket Client",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"url": {
					Title: "URL",
					Type:  "string",
				},
			},
			Required:             []string{"url"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ModuleConfig) (common.Module, error) {
			params := config.Params
			urlString, err := params.GetString("url")
			if err != nil {
				return nil, fmt.Errorf("websocket.client url error: %w", err)
			}

			parsedURL, err := url.Parse(urlString)
			if err != nil {
				return nil, fmt.Errorf("websocket.client url error: %w", err)
			}

			if parsedURL.Scheme != "ws" && parsedURL.Scheme != "wss" {
				return nil, fmt.Errorf("websocket.client url error: scheme must be ws or wss")
			}

			return &WebSocketClient{URL: *parsedURL, config: config, logger: CreateLogger(config)}, nil
		},
	})
}

func (wc *WebSocketClient) Id() string {
	return wc.config.Id
}

func (wc *WebSocketClient) Type() string {
	return wc.config.Type
}

func (wc *WebSocketClient) SetupConn() error {
	conn, _, err := websocket.DefaultDialer.Dial(wc.URL.String(), nil)
	wc.conn = conn
	return err
}

func (wc *WebSocketClient) Start(ctx context.Context, router common.RouteIO) error {
	wc.logger.Debug("running")
	wc.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	wc.ctx = moduleContext
	wc.cancel = cancel

	err := wc.SetupConn()
	if err != nil {
		return err
	}
	go wc.readLoop()

	<-wc.ctx.Done()
	wc.logger.Debug("done")
	if wc.conn != nil {
		wc.conn.Close()
	}
	return nil
}

func (wc *WebSocketClient) readLoop() {
	for {
		if wc.conn == nil {
			wc.SetupConn()
			continue
		}

		messageType, message, err := wc.conn.ReadMessage()
		if err != nil {
			wc.logger.Error("read error", "error", err)
			continue
		}
		if wc.router != nil {
			switch messageType {
			case websocket.TextMessage:
				wc.router.HandleInput(wc.ctx, wc.Id(), string(message))
			case websocket.BinaryMessage:
				wc.router.HandleInput(wc.ctx, wc.Id(), message)
			default:
				wc.logger.Warn("unsupported message type received", "messageType", messageType)
			}
		} else {
			wc.logger.Error("input received but no router is configured")
			continue
		}
	}
}

func (wc *WebSocketClient) outputBytes(ctx context.Context, payload []byte) error {
	if wc.conn == nil {
		return errors.New("websocket.client client is not setup")
	}
	err := wc.conn.WriteMessage(websocket.BinaryMessage, payload)

	if err != nil {
		return err
	}
	return nil
}

func (wc *WebSocketClient) outputString(ctx context.Context, payload string) error {
	if wc.conn == nil {
		return errors.New("websocket.client client is not setup")
	}
	err := wc.conn.WriteMessage(websocket.TextMessage, []byte(payload))

	if err != nil {
		return err
	}
	return nil
}

func (wc *WebSocketClient) Output(ctx context.Context, payload any) error {

	payloadBytes, ok := common.GetAnyAsByteSlice(payload)
	if ok {
		return wc.outputBytes(ctx, payloadBytes)
	} else {
		payloadString, ok := common.GetAnyAs[string](payload)
		if ok {
			return wc.outputString(ctx, payloadString)
		} else {
			return errors.New("websocket.client payload must be string or []byte")
		}
	}
}

func (wc *WebSocketClient) Stop() {
	wc.cancel()
}
