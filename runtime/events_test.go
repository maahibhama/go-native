package runtime

import "testing"

func TestEventRegistry(t *testing.T) {
	r := NewEventRegistry()
	calls := 0
	id := r.Register(func() { calls++ })
	if !r.Dispatch(id) || calls != 1 {
		t.Fatal("handler not dispatched")
	}
	r.Release(id)
	if r.Dispatch(id) {
		t.Fatal("released handler dispatched")
	}
}
