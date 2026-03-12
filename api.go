package showbridge

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
)

func (r *Router) handleConfigHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	configJSON, err := json.Marshal(r.runningConfig)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(configJSON)
}

//go:embed schema
var schema embed.FS

func (r *Router) handleSchemaHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	schemaName := req.PathValue("schema")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	configSchema, err := schema.ReadFile(fmt.Sprintf("schema/%s.schema.json", schemaName))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(configSchema)
}
