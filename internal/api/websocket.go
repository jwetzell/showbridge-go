package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/coder/websocket"
	"github.com/jwetzell/showbridge-go/internal/events"
)

type WebsocketEventDestination struct {
	conn *websocket.Conn
}

func (d WebsocketEventDestination) Send(ctx context.Context, event events.Event) error {
	eventJSON, err := event.ToJSON()
	if err != nil {
		return err
	}
	return d.conn.Write(ctx, websocket.MessageText, eventJSON)
}

func (d WebsocketEventDestination) Is(dest events.EventDestination) bool {
	other, ok := dest.(WebsocketEventDestination)
	if !ok {
		return false
	}
	return d.conn == other.conn
}

func (as *ApiServer) handleWebsocket(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	conn, err := websocket.Accept(w, req, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		as.logger.Error("websocket upgrade error", "error", err)
		return
	}
	defer conn.CloseNow()

	eventDestination := WebsocketEventDestination{conn: conn}

	as.eventRouter.AddEventDestination(eventDestination)
READ_LOOP:
	for {
		messageType, message, err := conn.Read(ctx)
		if err != nil {

			as.logger.Error("websocket read error", "error", err)
			break READ_LOOP
		}

		switch messageType {
		case websocket.MessageText, websocket.MessageBinary:
			event := events.Event{}
			err = json.Unmarshal(message, &event)
			if err != nil {
				as.logger.Error("websocket message unmarshal error", "error", err)
				continue
			}
			as.eventRouter.HandleEvent(event, WebsocketEventDestination{conn: conn})
		default:
			as.logger.Warn("unsupported websocket message type", "type", messageType)
			continue
		}

	}
	// NOTE(jwetzell): remove ws connection
	as.eventRouter.RemoveEventDestination(eventDestination)
}
