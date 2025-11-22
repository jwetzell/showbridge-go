package showbridge

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type HTTPServer struct {
	config ModuleConfig
	Port   uint16
	router *Router
}

type ResponseData struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "net.http.server",
		New: func(config ModuleConfig) (Module, error) {
			params := config.Params
			port, ok := params["port"]
			if !ok {
				return nil, fmt.Errorf("net.http.server requires a port parameter")
			}

			portNum, ok := port.(float64)

			if !ok {
				return nil, fmt.Errorf("net.http.server port must be uint16")
			}

			return &HTTPServer{Port: uint16(portNum), config: config}, nil
		},
	})
}

func (hs *HTTPServer) Id() string {
	return hs.config.Id
}

func (hs *HTTPServer) Type() string {
	return hs.config.Type
}

func (hs *HTTPServer) RegisterRouter(router *Router) {
	hs.router = router
}

func (hs *HTTPServer) HandleDefault(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := ResponseData{
		Message: "routing successful",
		Status:  "ok",
	}

	if hs.router != nil {
		routingErrors := hs.router.HandleInput(hs.config.Id, r)
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
		<-hs.router.Context.Done()
		slog.Debug("router context done in module", "id", hs.config.Id)
		httpServer.Close()
	}()

	err := httpServer.ListenAndServe()
	slog.Debug("net.http.server closed", "id", hs.config.Id)
	if err != nil {
		return err
	}

	<-hs.router.Context.Done()
	return nil
}

func (hs *HTTPServer) Output(payload any) error {
	return fmt.Errorf("net.http.server output is not implemented")
}
