package showbridge

import (
	"time"

	"github.com/jwetzell/showbridge-go/events"
)

func (r *Router) HandleEvent(event events.Event, sender events.EventDestination) {
	switch event.Type {
	case "ping":
		r.unicastEvent(events.Event{Type: "pong", Data: map[string]any{
			"timestamp": time.Now().UnixMilli(),
		}}, sender)
	default:
		r.logger.Warn("unknown event type", "eventType", event.Type)
	}
}

func (r *Router) AddEventDestination(dest events.EventDestination) {
	r.eventDestinationsMu.Lock()
	defer r.eventDestinationsMu.Unlock()
	r.eventDestinations = append(r.eventDestinations, dest)
}

func (r *Router) RemoveEventDestination(dest events.EventDestination) {
	r.eventDestinationsMu.Lock()
	defer r.eventDestinationsMu.Unlock()
	for i, d := range r.eventDestinations {
		if d.Is(dest) {
			r.eventDestinations = append(r.eventDestinations[:i], r.eventDestinations[i+1:]...)
			break
		}
	}
}

func (r *Router) unicastEvent(event events.Event, dest events.EventDestination) {
	err := dest.Send(event)
	if err != nil {
		r.logger.Error("failed to send event", "error", err)
	}
}

func (r *Router) broadcastEvent(event events.Event) {
	r.eventDestinationsMu.Lock()
	defer r.eventDestinationsMu.Unlock()
	for _, dest := range r.eventDestinations {
		err := dest.Send(event)
		if err != nil {
			r.logger.Error("failed to send event", "error", err)
		}
	}
}
