package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"sync"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/gorilla/websocket"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
)

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

type WebSocketClient struct {
	config       config.ModuleConfig
	URL          url.URL
	ctx          context.Context
	conn         *websocket.Conn
	inputHandler common.InputHandler
	logger       *slog.Logger
	cancel       context.CancelFunc
	connMu       sync.Mutex
}

func (wc *WebSocketClient) Id() string {
	return wc.config.Id
}

func (wc *WebSocketClient) Type() string {
	return wc.config.Type
}

func (wc *WebSocketClient) SetupConn() error {
	wc.connMu.Lock()
	defer wc.connMu.Unlock()
	if wc.conn != nil {
		wc.conn.Close()
	}
	conn, _, err := websocket.DefaultDialer.Dial(wc.URL.String(), nil)

	if err != nil {
		return fmt.Errorf("websocket.client dial error: %w", err)
	}

	conn.SetCloseHandler(func(code int, text string) error {
		// NOTE(jwetzell): attempt to send close message back to server before closing connection
		err := wc.conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			time.Now().Add(time.Minute),
		)
		wc.connMu.Lock()
		defer wc.connMu.Unlock()
		if wc.conn != nil {
			wc.conn.Close()
		}
		return err
	})

	wc.conn = conn
	return err
}

func (wc *WebSocketClient) Start(ctx context.Context, inputHandler common.InputHandler) error {
	wc.logger.Debug("running")
	wc.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	wc.ctx = moduleContext
	wc.cancel = cancel

	for wc.ctx.Err() == nil {
		err := wc.SetupConn()
		if err != nil {
			wc.logger.Error("connection error", "error", err)
		} else {
			// NOTE(jwetzell): enter read loop until an error occurs
			wc.readLoop()
		}
		// NOTE(jwetzell): if connection is lost or read error wait before trying again
		time.Sleep(2 * time.Second)
	}
	<-wc.ctx.Done()
	wc.logger.Debug("done")
	return nil
}

func (wc *WebSocketClient) readLoop() {
	for wc.ctx.Err() == nil {
		if wc.conn == nil {
			wc.logger.Error("websocket connection is not established")
			return
		}
		// TODO(jwetzell): other ways to timeout?
		wc.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		messageType, message, err := wc.conn.ReadMessage()
		if err != nil {
			return
		}
		if wc.inputHandler != nil {
			switch messageType {
			case websocket.CloseMessage:
				return
			case websocket.PingMessage:
				err := wc.conn.WriteMessage(websocket.PongMessage, nil)
				if err != nil {
					wc.logger.Error("websocket pong error", "error", err)
					return
				}
			case websocket.TextMessage:
				wc.inputHandler(wc.ctx, wc.Id(), string(message))
			case websocket.BinaryMessage:
				wc.inputHandler(wc.ctx, wc.Id(), message)
			default:
				wc.logger.Warn("unsupported message type received", "messageType", messageType)
			}
		} else {
			wc.logger.Error("input received but no input handler is configured")
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
	wc.connMu.Lock()
	defer wc.connMu.Unlock()
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
	if wc.cancel != nil {
		defer wc.cancel()
	}
	wc.connMu.Lock()
	defer wc.connMu.Unlock()
	if wc.conn != nil {
		err := wc.conn.WriteControl(websocket.CloseMessage, nil, time.Now().Add(time.Minute))
		if err != nil {
			wc.logger.Error("websocket close error", "error", err)
		}
		wc.conn.Close()
	}
}
