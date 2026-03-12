package showbridge

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (r *Router) handleWebsocket(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		r.logger.Error("websocket upgrade error", "error", err)
		return
	}
	defer conn.Close()

	r.wsConnsMu.Lock()
	r.wsConns = append(r.wsConns, conn)
	r.wsConnsMu.Unlock()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			r.logger.Error("websocket read error", "error", err)
			break
		}

		event := Event{}
		err = json.Unmarshal(message, &event)
		if err != nil {
			r.logger.Error("websocket message unmarshal error", "error", err)
			continue
		}
		r.handleEvent(event)
	}
	r.wsConnsMu.Lock()
	for i, c := range r.wsConns {
		if c == conn {
			r.wsConns = append(r.wsConns[:i], r.wsConns[i+1:]...)
			break
		}
	}
	r.wsConnsMu.Unlock()
}
