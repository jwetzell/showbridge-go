package showbridge

import (
	"time"

	"github.com/jwetzell/showbridge-go/internal/common"
)

func (r *Router) HandleEvent(event common.Event, sender common.EventDestination) {
	switch event.Type {
	case "ping":
		r.unicastEvent(common.Event{Type: "pong", Data: map[string]any{
			"timestamp": time.Now().UnixMilli(),
		}}, sender)
	default:
		r.logger.Warn("unknown event type", "eventType", event.Type)
	}
}

func (r *Router) AddEventDestination(dest common.EventDestination) {
	r.eventDestinationsMu.Lock()
	defer r.eventDestinationsMu.Unlock()
	r.eventDestinations = append(r.eventDestinations, dest)
}

func (r *Router) RemoveEventDestination(dest common.EventDestination) {
	r.eventDestinationsMu.Lock()
	defer r.eventDestinationsMu.Unlock()
	for i, d := range r.eventDestinations {
		if d.Is(dest) {
			r.eventDestinations = append(r.eventDestinations[:i], r.eventDestinations[i+1:]...)
			break
		}
	}
}

func (r *Router) unicastEvent(event common.Event, dest common.EventDestination) {
	err := dest.Send(event)
	if err != nil {
		r.logger.Error("failed to send event", "error", err)
	}
}

func (r *Router) broadcastEvent(event common.Event) {
	r.eventDestinationsMu.Lock()
	defer r.eventDestinationsMu.Unlock()
	for _, dest := range r.eventDestinations {
		err := dest.Send(event)
		if err != nil {
			r.logger.Error("failed to send event", "error", err)
		}
	}
}
