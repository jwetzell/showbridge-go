package module

import (
	"context"
	"fmt"
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
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.http.client",
		New: func(ctx context.Context, config config.ModuleConfig, router route.RouteIO) (Module, error) {

			return &HTTPClient{config: config, ctx: ctx, router: router}, nil
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
	slog.Debug("router context done in module", "id", hc.Id())
	return nil
}

func (hc *HTTPClient) Output(payload any) error {

	payloadRequest, ok := payload.(*http.Request)

	if !ok {
		return fmt.Errorf("net.http.client is only able to output an http.Request")
	}

	if hc.client == nil {
		return fmt.Errorf("net.http.client client is nil")
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
