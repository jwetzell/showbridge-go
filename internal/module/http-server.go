package module

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/route"
)

type HTTPServer struct {
	config config.ModuleConfig
	Port   uint16
	ctx    context.Context
	router route.RouteIO
	logger *slog.Logger
}

type ResponseData struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "http.server",
		New: func(ctx context.Context, config config.ModuleConfig, router route.RouteIO) (Module, error) {
			params := config.Params
			port, ok := params["port"]
			if !ok {
				return nil, errors.New("http.server requires a port parameter")
			}

			portNum, ok := port.(float64)

			if !ok {
				return nil, errors.New("http.server port must be uint16")
			}

			return &HTTPServer{Port: uint16(portNum), config: config, ctx: ctx, router: router, logger: slog.Default().With("component", "module", "id", config.Id)}, nil
		},
	})
}

func (hs *HTTPServer) Id() string {
	return hs.config.Id
}

func (hs *HTTPServer) Type() string {
	return hs.config.Type
}

func (hs *HTTPServer) HandleDefault(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := ResponseData{
		Message: "routing successful",
		Status:  "ok",
	}

	if hs.router != nil {
		routingErrors := hs.router.HandleInput(hs.Id(), r)
		if routingErrors != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.Status = "error"
			response.Message = "routing failed"
		} else {
			w.WriteHeader(http.StatusOK)
			response.Message = "routing successful"
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = "no router registered"
		response.Status = "error"
	}

	json.NewEncoder(w).Encode(response)
}

func (hs *HTTPServer) Run() error {
	http.HandleFunc("/", hs.HandleDefault)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", hs.Port),
		Handler: http.DefaultServeMux,
	}

	go func() {
		<-hs.ctx.Done()
		hs.logger.Debug("router context done in module")
		httpServer.Close()
	}()

	err := httpServer.ListenAndServe()
	// TODO(jwetzell): handle server closed error differently
	if err != nil {
		return err
	}

	<-hs.ctx.Done()
	return nil
}

func (hs *HTTPServer) Output(payload any) error {
	return errors.New("http.server output is not implemented")
}
