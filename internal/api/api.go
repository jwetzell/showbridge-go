package api

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/jwetzell/showbridge-go/internal/schema"
)

type ApiServer struct {
	config             config.ApiConfig
	serverMu           sync.Mutex
	server             *http.Server
	shutdown           context.CancelFunc
	logger             *slog.Logger
	configurableRouter config.Configurable
	eventRouter        common.EventRouter
}

func NewApiServer(configurableRouter config.Configurable, eventRouter common.EventRouter) *ApiServer {
	return &ApiServer{
		configurableRouter: configurableRouter,
		eventRouter:        eventRouter,
		logger:             slog.Default().With("component", "api"),
	}
}

func (as *ApiServer) Start(config config.ApiConfig) {
	as.config = config
	if !as.config.Enabled {
		as.logger.Warn("not enabled")
		return
	}
	as.logger.Debug("starting", "port", as.config.Port)
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", as.handleWebsocket)
	mux.HandleFunc("/health", as.handleHealthHTTP)
	mux.HandleFunc("/api/v1/config", as.handleConfigHTTP)
	mux.HandleFunc("/schema/config.schema.json", handleConfigSchema)
	mux.HandleFunc("/schema/routes.schema.json", handleRoutesSchema)
	mux.HandleFunc("/schema/modules.schema.json", handleModulesSchema)
	mux.HandleFunc("/schema/processors.schema.json", handleProcessorsSchema)

	as.serverMu.Lock()
	defer as.serverMu.Unlock()
	as.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", as.config.Port),
		Handler: mux,
	}

	go func() {
		as.server.ListenAndServe()
		as.shutdown()
	}()
}

func (as *ApiServer) Stop() {
	if as.server == nil {
		return
	}
	as.logger.Debug("stopping")
	as.serverMu.Lock()
	defer as.serverMu.Unlock()
	if as.server != nil {
		apiShutdownCtx, apiShutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		as.shutdown = apiShutdownCancel
		as.server.Shutdown(apiShutdownCtx)
		<-apiShutdownCtx.Done()
		as.server = nil
	}
}

func (as *ApiServer) handleHealthHTTP(w http.ResponseWriter, req *http.Request) {
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

func (as *ApiServer) handleConfigHTTP(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodGet:
		configJSON, err := json.Marshal(as.configurableRouter.GetRunningConfig())
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Write(configJSON)
	case http.MethodPut:
		//TODO(jwetzell): again way too much marshaling
		cfgBytes, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		cfgMap := make(map[string]any)
		err = json.Unmarshal(cfgBytes, &cfgMap)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		err = schema.ApplyDefaults(&cfgMap)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = schema.ValidateConfig(cfgMap)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		validCfgBytes, err := json.Marshal(cfgMap)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		var newConfig config.Config
		err = json.Unmarshal(validCfgBytes, &newConfig)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		err, moduleErrors, routeErrors := as.configurableRouter.UpdateConfig(newConfig, true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		if len(moduleErrors) > 0 || len(routeErrors) > 0 {
			errorResponse := struct {
				ModuleErrors []config.ModuleError `json:"moduleErrors,omitempty"`
				RouteErrors  []config.RouteError  `json:"routeErrors,omitempty"`
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

func handleConfigSchema(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		schemaJSON, err := json.Marshal(schema.ConfigSchema)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Write(schemaJSON)
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

func handleRoutesSchema(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		schemaJSON, err := json.Marshal(schema.RoutesConfigSchema)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Write(schemaJSON)
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

func handleModulesSchema(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		schemaJSON, err := json.Marshal(schema.GetModulesSchema())
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Write(schemaJSON)
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

func handleProcessorsSchema(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		schemaJSON, err := json.Marshal(schema.GetProcessorsSchema())
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Write(schemaJSON)
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
