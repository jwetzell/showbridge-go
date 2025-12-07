package showbridge

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jwetzell/showbridge-go/internal/config"
)

type HTTPClient struct {
	config config.ModuleConfig
	router *Router
	client *http.Client
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.http.client",
		New: func(config config.ModuleConfig, router *Router) (Module, error) {

			return &HTTPClient{config: config, router: router}, nil
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

	<-hc.router.Context.Done()
	slog.Debug("router context done in module", "id", hc.config.Id)
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
		hc.router.HandleInput(hc.config.Id, response)
	}

	return nil
}
