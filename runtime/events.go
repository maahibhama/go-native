package runtime

import (
	"github.com/go-native/go-native/ui"
	"sync"
	"sync/atomic"
)

// EventRegistry owns Go callbacks referenced by stable integer IDs.
type EventRegistry struct {
	next            atomic.Uint64
	mu              sync.RWMutex
	handlers        map[ui.HandlerID]func()
	valueHandlers   map[ui.HandlerID]func(string)
	boolHandlers    map[ui.HandlerID]func(bool)
	gestureHandlers map[ui.HandlerID]func(ui.GestureEvent)
}

func NewEventRegistry() *EventRegistry {
	return &EventRegistry{handlers: make(map[ui.HandlerID]func()), valueHandlers: make(map[ui.HandlerID]func(string)), boolHandlers: make(map[ui.HandlerID]func(bool)), gestureHandlers: make(map[ui.HandlerID]func(ui.GestureEvent))}
}
func (r *EventRegistry) RegisterValue(fn func(string)) ui.HandlerID {
	if fn == nil {
		return 0
	}
	id := ui.HandlerID(r.next.Add(1))
	r.mu.Lock()
	r.valueHandlers[id] = fn
	r.mu.Unlock()
	return id
}
func (r *EventRegistry) DispatchValue(id ui.HandlerID, value string) bool {
	r.mu.RLock()
	fn := r.valueHandlers[id]
	r.mu.RUnlock()
	if fn == nil {
		return false
	}
	fn(value)
	return true
}
func (r *EventRegistry) ReplaceValue(id ui.HandlerID, fn func(string)) {
	r.mu.Lock()
	r.valueHandlers[id] = fn
	r.mu.Unlock()
}
func (r *EventRegistry) RegisterBool(fn func(bool)) ui.HandlerID {
	if fn == nil {
		return 0
	}
	id := ui.HandlerID(r.next.Add(1))
	r.mu.Lock()
	r.boolHandlers[id] = fn
	r.mu.Unlock()
	return id
}
func (r *EventRegistry) DispatchBool(id ui.HandlerID, value bool) bool {
	r.mu.RLock()
	fn := r.boolHandlers[id]
	r.mu.RUnlock()
	if fn == nil {
		return false
	}
	fn(value)
	return true
}
func (r *EventRegistry) ReplaceBool(id ui.HandlerID, fn func(bool)) {
	r.mu.Lock()
	r.boolHandlers[id] = fn
	r.mu.Unlock()
}
func (r *EventRegistry) RegisterGesture(fn func(ui.GestureEvent)) ui.HandlerID {
	if fn == nil {
		return 0
	}
	id := ui.HandlerID(r.next.Add(1))
	r.mu.Lock()
	r.gestureHandlers[id] = fn
	r.mu.Unlock()
	return id
}
func (r *EventRegistry) DispatchGesture(id ui.HandlerID, event ui.GestureEvent) bool {
	r.mu.RLock()
	fn := r.gestureHandlers[id]
	r.mu.RUnlock()
	if fn == nil {
		return false
	}
	fn(event)
	return true
}
func (r *EventRegistry) ReplaceGesture(id ui.HandlerID, fn func(ui.GestureEvent)) {
	r.mu.Lock()
	r.gestureHandlers[id] = fn
	r.mu.Unlock()
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
	delete(r.valueHandlers, id)
	delete(r.boolHandlers, id)
	delete(r.gestureHandlers, id)
	r.mu.Unlock()
}
