package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
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
	conn, _, err := websocket.DefaultDialer.Dial(wc.URL.String(), nil)
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
			wc.logger.Debug("websocket connection established entering read loop")
			wc.readLoop()
		}
		// NOTE(jwetzell): if connection is lost or read error wait before trying again
		time.Sleep(2 * time.Second)
	}
	<-wc.ctx.Done()
	return nil
}

func (wc *WebSocketClient) readLoop() {
	for wc.ctx.Err() == nil {
		if wc.conn == nil {
			wc.logger.Error("websocket connection is not established")
			return
		}
		wc.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		messageType, message, err := wc.conn.ReadMessage()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok {
				// NOTE(jwetzell) we hit deadline
				if opErr.Timeout() {
					continue
				}
				// NOTE(jwetzell) connection was closed
				if errors.Is(opErr, net.ErrClosed) {
					continue
				}
			}
			wc.logger.Error("websocket read error", "error", err)
			return
		}
		if wc.inputHandler != nil {
			switch messageType {
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
		wc.cancel()
	}
	wc.connMu.Lock()
	defer wc.connMu.Unlock()
	if wc.conn != nil {
		wc.conn.Close()
		wc.conn = nil
	}
	wc.logger.Debug("done")
}
