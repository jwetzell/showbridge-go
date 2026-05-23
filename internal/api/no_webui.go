//go:build !showbridge_webui

package api

import (
	"net/http"
)

func handleWebUI(w http.ResponseWriter, req *http.Request) {

	http.Error(w, "Web UI is not enabled", http.StatusNotImplemented)
}
