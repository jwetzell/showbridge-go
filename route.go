package showbridge

import "log/slog"

type Route struct {
	index  int
	Input  string
	Output string
	router *Router
}

type RouteConfig struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

func NewRoute(index int, config RouteConfig, router *Router) *Route {
	return &Route{Input: config.Input, Output: config.Output, router: router, index: index}
}

func (r *Route) HandleInput(sourceId string, payload any) {
	slog.Debug("route input", "index", r.index, "source", sourceId, "payload", payload)
	r.HandleOutput(payload)
}

func (r *Route) HandleOutput(payload any) {
	slog.Debug("route output", "index", r.index, "destination", r.Output, "payload", payload)
	err := r.router.HandleOutput(r.Output, payload)
	if err != nil {
		slog.Error("problem with route output", "error", err.Error())
	}
}
