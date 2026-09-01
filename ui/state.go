package ui

import "sync"

// Scheduler coalesces state changes into rendering passes.
type Scheduler interface{ Schedule() }

var scheduler struct {
	sync.RWMutex
	value Scheduler
}

// SetScheduler connects state changes to an application runtime.
func SetScheduler(s Scheduler) { scheduler.Lock(); scheduler.value = s; scheduler.Unlock() }

func schedule() {
	scheduler.RLock()
	s := scheduler.value
	scheduler.RUnlock()
	if s != nil {
		s.Schedule()
	}
}

// State is a thread-safe observable value.
type State[T any] struct {
	mu    sync.RWMutex
	value T
}

// NewState creates state with an initial value.
func NewState[T any](value T) *State[T] { return &State[T]{value: value} }

// Get returns the current value.
func (s *State[T]) Get() T { s.mu.RLock(); defer s.mu.RUnlock(); return s.value }

// Set replaces the value and schedules a rendering pass.
func (s *State[T]) Set(value T) { s.mu.Lock(); s.value = value; s.mu.Unlock(); schedule() }

// Update atomically replaces the value using update and schedules a rendering pass.
func (s *State[T]) Update(update func(T) T) {
	s.mu.Lock()
	s.value = update(s.value)
	s.mu.Unlock()
	schedule()
}
