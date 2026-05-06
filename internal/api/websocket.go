package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jwetzell/showbridge-go/internal/common"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketEventDestination struct {
	conn *websocket.Conn
}

func (d WebsocketEventDestination) Send(event common.Event) error {
	eventJSON, err := event.ToJSON()
	if err != nil {
		return err
	}
	return d.conn.WriteMessage(websocket.TextMessage, eventJSON)
}

func (d WebsocketEventDestination) Is(dest common.EventDestination) bool {
	other, ok := dest.(WebsocketEventDestination)
	if !ok {
		return false
	}
	return d.conn == other.conn
}

func (as *ApiServer) handleWebsocket(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		as.logger.Error("websocket upgrade error", "error", err)
		return
	}
	defer conn.Close()

	eventDestination := WebsocketEventDestination{conn: conn}

	as.eventRouter.AddEventDestination(eventDestination)
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
			event := common.Event{}
			err = json.Unmarshal(message, &event)
			if err != nil {
				as.logger.Error("websocket message unmarshal error", "error", err)
				continue
			}
			as.eventRouter.HandleEvent(event, WebsocketEventDestination{conn: conn})
		case websocket.CloseMessage:
			break READ_LOOP
		case websocket.PingMessage:
			err = conn.WriteMessage(websocket.PongMessage, nil)
			if err != nil {
				as.logger.Error("websocket pong error", "error", err)
			}
		default:
			as.logger.Warn("unsupported websocket message type", "type", messageType)
			continue
		}

	}
	//NOTE(jwetzell): remove ws connection
	as.eventRouter.RemoveEventDestination(eventDestination)
}
