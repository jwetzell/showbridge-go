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
READ_LOOP:
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			_, ok := err.(*websocket.CloseError)
			if ok {
				break READ_LOOP
			}
		}

		switch messageType {
		case websocket.TextMessage, websocket.BinaryMessage:
			event := Event{}
			err = json.Unmarshal(message, &event)
			if err != nil {
				r.logger.Error("websocket message unmarshal error", "error", err)
				continue
			}
			r.handleEvent(event, conn)
		case websocket.CloseMessage:
			break READ_LOOP
		case websocket.PingMessage:
			err = conn.WriteMessage(websocket.PongMessage, nil)
			if err != nil {
				r.logger.Error("websocket pong error", "error", err)
			}
		default:
			r.logger.Warn("unsupported websocket message type", "type", messageType)
			continue
		}

	}
	//NOTE(jwetzell): remove ws connection
	r.wsConnsMu.Lock()
	for i, c := range r.wsConns {
		if c == conn {
			r.wsConns = append(r.wsConns[:i], r.wsConns[i+1:]...)
			break
		}
	}
	r.wsConnsMu.Unlock()
}
