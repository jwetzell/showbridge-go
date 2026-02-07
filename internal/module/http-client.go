package module

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/route"
)

type HTTPClient struct {
	config config.ModuleConfig
	ctx    context.Context
	client *http.Client
	router route.RouteIO
	logger *slog.Logger
	cancel context.CancelFunc
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "http.client",
		New: func(config config.ModuleConfig) (Module, error) {

			return &HTTPClient{config: config, logger: CreateLogger(config)}, nil
		},
	})
}

func (hc *HTTPClient) Id() string {
	return hc.config.Id
}

func (hc *HTTPClient) Type() string {
	return hc.config.Type
}

func (hc *HTTPClient) Run(ctx context.Context) error {
	router, ok := ctx.Value(route.RouterContextKey).(route.RouteIO)

	if !ok {
		return errors.New("http.client unable to get router from context")
	}
	hc.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	hc.ctx = moduleContext
	hc.cancel = cancel

	hc.client = &http.Client{
		Timeout: 10 * time.Second,
	}

	<-hc.ctx.Done()
	hc.logger.Debug("done")
	return nil
}

func (hc *HTTPClient) Output(ctx context.Context, payload any) error {

	payloadRequest, ok := payload.(*http.Request)

	if !ok {
		return errors.New("http.client is only able to output an http.Request")
	}

	if hc.client == nil {
		return errors.New("http.client client is nil")
	}

	response, err := hc.client.Do(payloadRequest)

	if err != nil {
		return err
	}

	if hc.router != nil {
		hc.router.HandleInput(hc.ctx, hc.Id(), response)
	}

	return nil
}

func (hc *HTTPClient) Stop() {
	hc.cancel()
}
