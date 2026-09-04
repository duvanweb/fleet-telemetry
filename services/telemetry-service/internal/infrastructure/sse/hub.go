package sse

import (
	"sync"

	"github.com/oklog/ulid/v2"
)

// Event is a single SSE message with an event type and JSON payload.
type Event struct {
	Type    string
	Payload []byte
}

// Hub maintains the set of active SSE subscribers and broadcasts events to them.
type Hub struct {
	mu      sync.RWMutex
	clients map[string]chan Event
}

// NewHub creates and returns a new Hub.
func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]chan Event),
	}
}

// Broadcast sends an event to all current subscribers. Slow clients whose
// channel buffer is full are skipped (non-blocking).
func (h *Hub) Broadcast(event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, ch := range h.clients {
		select {
		case ch <- event:
		default:
		}
	}
}

// Subscribe registers a new subscriber and returns its ID and receive channel.
func (h *Hub) Subscribe() (string, <-chan Event) {
	id := ulid.Make().String()
	ch := make(chan Event, 64)

	h.mu.Lock()
	h.clients[id] = ch
	h.mu.Unlock()

	return id, ch
}

// Unsubscribe removes a subscriber by ID and closes its channel.
func (h *Hub) Unsubscribe(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if ch, ok := h.clients[id]; ok {
		close(ch)
		delete(h.clients, id)
	}
}
