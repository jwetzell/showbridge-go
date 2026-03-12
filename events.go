package showbridge

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Event struct {
	Type  string `json:"type"`
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func (e Event) toJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (r *Router) handleEvent(event Event, sender *websocket.Conn) {
	switch event.Type {
	case "ping":
		r.unicastEvent(Event{Type: "pong"}, sender)
	default:
		r.logger.Warn("unknown event type", "eventType", event.Type)
	}
}

func (r *Router) unicastEvent(event Event, conn *websocket.Conn) {
	eventJSON, err := event.toJSON()
	if err != nil {
		r.logger.Error("failed to marshal event to JSON", "error", err)
		return
	}
	err = conn.WriteMessage(websocket.TextMessage, eventJSON)
	if err != nil {
		r.logger.Error("failed to write message to websocket connection", "error", err)
	}
}

func (r *Router) broadcastEvent(event Event, excluded ...*websocket.Conn) {
	eventJSON, err := event.toJSON()
	if err != nil {
		r.logger.Error("failed to marshal event to JSON", "error", err)
		return
	}
	r.wsConnsMu.Lock()
	defer r.wsConnsMu.Unlock()
	for _, conn := range r.wsConns {
		exclude := false
		for _, excludedConn := range excluded {
			if conn == excludedConn {
				exclude = true
				break
			}
		}
		if exclude {
			continue
		}
		err := conn.WriteMessage(websocket.TextMessage, eventJSON)
		if err != nil {
			r.logger.Error("failed to write message to websocket connection", "error", err)
		}
	}
}
