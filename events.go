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

func (r *Router) handleEvent(event Event) {
	switch event.Type {
	case "ping":
		r.broadcastEvent(Event{Type: "pong"})
	default:
		r.logger.Warn("unknown event type", "eventType", event.Type)
	}
}

func (r *Router) broadcastEvent(event Event) {
	eventJSON, err := event.toJSON()
	if err != nil {
		r.logger.Error("failed to marshal event to JSON", "error", err)
		return
	}
	r.wsConnsMu.Lock()
	defer r.wsConnsMu.Unlock()
	for _, conn := range r.wsConns {
		err := conn.WriteMessage(websocket.TextMessage, eventJSON)
		if err != nil {
			r.logger.Error("failed to write message to websocket connection", "error", err)
		}
	}
}
