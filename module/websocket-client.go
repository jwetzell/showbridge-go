package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/config"
	"github.com/jwetzell/showbridge-go/params"
	"github.com/jwetzell/showbridge-go/routing"
)

func init() {
	RegisterModule(ModuleRegistration{
		Type:        "websocket.client",
		Title:       "WebSocket Client",
		Description: "Opens a WebSocket connection to the specified server",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"url": {
					Title:       "URL",
					Description: "the URL of the WebSocket server to connect",
					Type:        "string",
				},
			},
			Required:             []string{"url"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ModuleConfig) (routing.Module, error) {
			moduleParams := config.Params
			urlString, err := moduleParams.GetString("url")
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
	config config.ModuleConfig
	URL    url.URL
	ctx    context.Context
	conn   *websocket.Conn
	router routing.Router
	logger *slog.Logger
	cancel context.CancelFunc
	connMu sync.Mutex
}

func (wc *WebSocketClient) SetupConn() error {
	wc.connMu.Lock()
	defer wc.connMu.Unlock()
	if wc.conn != nil {
		err := wc.conn.CloseNow()
		if err != nil {
			wc.logger.Error("websocket client close error", "error", err)
		}
	}
	conn, _, err := websocket.Dial(wc.ctx, wc.URL.String(), nil)

	if err != nil {
		return fmt.Errorf("websocket.client dial error: %w", err)
	}

	wc.conn = conn
	return err
}

func (wc *WebSocketClient) Start(ctx context.Context, router routing.Router) error {
	wc.logger.Debug("running")
	wc.router = router
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
		messageType, message, err := wc.conn.Read(wc.ctx)
		if err != nil {
			wc.logger.Error("websocket read error", "error", err)
			return
		}
		if wc.router != nil {
			switch messageType {
			case websocket.MessageText:
				wc.router.HandleInput(wc.ctx, wc.config.Id, string(message))
			case websocket.MessageBinary:
				wc.router.HandleInput(wc.ctx, wc.config.Id, message)
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
	err := wc.conn.Write(wc.ctx, websocket.MessageBinary, payload)

	if err != nil {
		return err
	}
	return nil
}

func (wc *WebSocketClient) outputString(ctx context.Context, payload string) error {
	if wc.conn == nil {
		return errors.New("websocket.client client is not setup")
	}

	err := wc.conn.Write(wc.ctx, websocket.MessageText, []byte(payload))

	if err != nil {
		return err
	}
	return nil
}

func (wc *WebSocketClient) Output(ctx context.Context, payload any) error {
	wc.connMu.Lock()
	defer wc.connMu.Unlock()
	payloadString, ok := payload.(string)
	if ok {
		return wc.outputString(ctx, payloadString)
	}
	payloadBytes, ok := params.GetAnyAsByteSlice(payload)
	if ok {
		return wc.outputBytes(ctx, payloadBytes)
	}
	return errors.New("websocket.client payload must be string or []byte")
}

func (wc *WebSocketClient) Stop() {
	if wc.cancel != nil {
		defer wc.cancel()
	}
	wc.connMu.Lock()
	defer wc.connMu.Unlock()
	if wc.conn != nil {
		err := wc.conn.Close(websocket.StatusNormalClosure, "")
		if err != nil {
			wc.logger.Error("websocket close error", "error", err)
		}
		err = wc.conn.CloseNow()
		if err != nil {
			wc.logger.Error("websocket close error", "error", err)
		}
	}
}
