package common

import (
	"encoding/json"
)

type Event struct {
	Type  string `json:"type"`
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func (e Event) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

type EventDestination interface {
	Send(event Event) error
	Is(dest EventDestination) bool
}

type EventRouter interface {
	HandleEvent(event Event, source EventDestination)
	AddEventDestination(dest EventDestination)
	RemoveEventDestination(dest EventDestination)
}
