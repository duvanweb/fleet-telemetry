package controllers

import (
	"fmt"
	"net/http"

	"fleet/telemetry-service/internal/infrastructure/sse"
	"fleet/shared/pkg/logger"
)

// Events is the HTTP controller for the SSE event stream endpoint.
type Events struct {
	logger logger.Logger
	hub    *sse.Hub
}

// @Router /events [get]
// @Tags events
// @Summary Subscribe to the real-time fleet event stream (SSE).
// @Success 200 "Event stream opened."
// GetEvents streams real-time fleet events to the client using Server-Sent Events.
func (c *Events) GetEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		c.logger.Errorw(r.Context(), "streaming not supported by response writer")
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	id, ch := c.hub.Subscribe()
	defer c.hub.Unsubscribe(id)

	for {
		select {
		case <-r.Context().Done():
			return
		case event, ok := <-ch:
			if !ok {
				return
			}
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event.Type, event.Payload)
			flusher.Flush()
		}
	}
}

// NewEvents creates and returns a new Events controller.
func NewEvents(log logger.Logger, hub *sse.Hub) *Events {
	return &Events{logger: log, hub: hub}
}
