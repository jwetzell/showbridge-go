//go:build showbridge_webui

package api

import (
	"embed"
	"errors"
	"io/fs"
	"net/http"
	"path"
	"strings"

	_ "embed"
)

//go:embed webui
var webUIFS embed.FS

var fsPrefix = "webui/browser"

var index = path.Join(fsPrefix, "index.html")

func handleWebUI(w http.ResponseWriter, req *http.Request) {
	requestedPath := strings.TrimLeft(req.URL.Path, "/ui")
	switch req.Method {
	case http.MethodGet:
		var pathToLoad string
		stat, err := fs.Stat(webUIFS, path.Join(fsPrefix, requestedPath))
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				pathToLoad = index
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
		if stat != nil {
			if stat.IsDir() {
				pathToLoad = index
			} else {
				pathToLoad = path.Join(fsPrefix, requestedPath)
			}
		}
		file, err := webUIFS.ReadFile(pathToLoad)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")

		switch path.Ext(pathToLoad) {
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".html":
			w.Header().Set("Content-Type", "text/html")
		case ".ico":
			w.Header().Set("Content-Type", "image/x-icon")
		}
		w.Write(file)
		return

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
