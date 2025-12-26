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
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "http.client",
		New: func(ctx context.Context, config config.ModuleConfig, router route.RouteIO) (Module, error) {

			return &HTTPClient{config: config, ctx: ctx, router: router, logger: CreateLogger(config)}, nil
		},
	})
}

func (hc *HTTPClient) Id() string {
	return hc.config.Id
}

func (hc *HTTPClient) Type() string {
	return hc.config.Type
}

func (hc *HTTPClient) Run() error {

	hc.client = &http.Client{
		Timeout: 10 * time.Second,
	}

	<-hc.ctx.Done()
	hc.logger.Debug("router context done in module")
	return nil
}

func (hc *HTTPClient) Output(payload any) error {

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
		hc.router.HandleInput(hc.Id(), response)
	}

	return nil
}
