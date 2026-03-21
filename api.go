package showbridge

import (
	"context"
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/module"
	"github.com/jwetzell/showbridge-go/internal/route"
)

//go:embed schema
var schema embed.FS

func (r *Router) startAPIServer(config config.ApiConfig) {
	if !config.Enabled {
		r.logger.Warn("API not enabled")
		return
	}
	r.logger.Debug("starting API server", "port", config.Port)
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", r.handleWebsocket)
	mux.HandleFunc("/health", r.handleHealthHTTP)
	mux.Handle("/schema/{schema}", HandleFS(schema))
	mux.HandleFunc("/api/v1/config", r.handleConfigHTTP)

	r.apiServerMu.Lock()
	defer r.apiServerMu.Unlock()
	r.apiServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: mux,
	}

	go func() {
		r.apiServer.ListenAndServe()
		r.apiServerShutdown()
	}()
}

func (r *Router) stopAPIServer() {
	if r.apiServer == nil {
		return
	}
	r.logger.Debug("stopping API server")
	r.apiServerMu.Lock()
	defer r.apiServerMu.Unlock()
	if r.apiServer != nil {
		apiShutdownCtx, apiShutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		r.apiServerShutdown = apiShutdownCancel
		r.apiServer.Shutdown(apiShutdownCtx)
		<-apiShutdownCtx.Done()
		r.apiServer = nil
	}
}

func (r *Router) handleHealthHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
	default:
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleConfigHTTP(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodGet:
		configJSON, err := json.Marshal(r.runningConfig)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Write(configJSON)
	case http.MethodPut:
		if r.updatingConfig {
			http.Error(w, "Config update in progress.", http.StatusConflict)
			return
		}
		var newConfig config.Config
		err := json.NewDecoder(req.Body).Decode(&newConfig)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		moduleErrors, routeErrors := r.UpdateConfig(newConfig)
		if len(moduleErrors) > 0 || len(routeErrors) > 0 {
			errorResponse := struct {
				ModuleErrors []module.ModuleError `json:"moduleErrors,omitempty"`
				RouteErrors  []route.RouteError   `json:"routeErrors,omitempty"`
			}{
				ModuleErrors: moduleErrors,
				RouteErrors:  routeErrors,
			}
			errorResponseJSON, err := json.Marshal(errorResponse)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errorResponseJSON)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		r.ConfigChange <- newConfig
	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
	default:
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func HandleFS(fs fs.FS) http.HandlerFunc {
	handler := http.FileServerFS(fs)
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(w, req)
	}
}
