package runtime

import (
	"github.com/go-native/go-native/ui"
	"sync"
	"sync/atomic"
)

// EventRegistry owns Go callbacks referenced by stable integer IDs.
type EventRegistry struct {
	next     atomic.Uint64
	mu       sync.RWMutex
	handlers map[ui.HandlerID]func()
}

func NewEventRegistry() *EventRegistry {
	return &EventRegistry{handlers: make(map[ui.HandlerID]func())}
}
func (r *EventRegistry) Register(fn func()) ui.HandlerID {
	if fn == nil {
		return 0
	}
	id := ui.HandlerID(r.next.Add(1))
	r.mu.Lock()
	r.handlers[id] = fn
	r.mu.Unlock()
	return id
}
func (r *EventRegistry) Dispatch(id ui.HandlerID) bool {
	r.mu.RLock()
	fn := r.handlers[id]
	r.mu.RUnlock()
	if fn == nil {
		return false
	}
	fn()
	return true
}
func (r *EventRegistry) Replace(id ui.HandlerID, fn func()) {
	r.mu.Lock()
	r.handlers[id] = fn
	r.mu.Unlock()
}
func (r *EventRegistry) Release(id ui.HandlerID) {
	if id == 0 {
		return
	}
	r.mu.Lock()
	delete(r.handlers, id)
	r.mu.Unlock()
}
