package ui

import "sync"

type callbackStore struct {
	mu     sync.RWMutex
	values map[NodeID]func()
}

func (s *callbackStore) Store(id NodeID, fn func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.values == nil {
		s.values = make(map[NodeID]func())
	}
	s.values[id] = fn
}

func (s *callbackStore) Load(id NodeID) func() {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn := s.values[id]
	delete(s.values, id)
	return fn
}
