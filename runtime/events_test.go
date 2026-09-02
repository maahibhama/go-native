package runtime

import (
	"github.com/go-native/go-native/ui"
	"testing"
)

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

func TestValueEventRegistry(t *testing.T) {
	r := NewEventRegistry()
	got := ""
	id := r.RegisterValue(func(value string) { got = value })
	if !r.DispatchValue(id, "hello") || got != "hello" {
		t.Fatalf("value handler got %q", got)
	}
	r.ReplaceValue(id, func(value string) { got = "new:" + value })
	r.DispatchValue(id, "value")
	if got != "new:value" {
		t.Fatalf("replacement got %q", got)
	}
	r.Release(id)
	if r.DispatchValue(id, "ignored") {
		t.Fatal("released value handler dispatched")
	}
}

func TestBoolEventRegistry(t *testing.T) {
	r := NewEventRegistry()
	got := false
	id := r.RegisterBool(func(value bool) { got = value })
	if !r.DispatchBool(id, true) || !got {
		t.Fatal("bool handler not dispatched")
	}
	r.ReplaceBool(id, func(value bool) { got = !value })
	r.DispatchBool(id, true)
	if got {
		t.Fatal("bool replacement not used")
	}
	r.Release(id)
	if r.DispatchBool(id, true) {
		t.Fatal("released bool handler dispatched")
	}
}

func TestGestureEventRegistryPreservesFourValues(t *testing.T) {
	r := NewEventRegistry()
	var got ui.GestureEvent
	id := r.RegisterGesture(func(event ui.GestureEvent) { got = event })
	want := ui.GestureEvent{TranslationX: 1, TranslationY: 2, VelocityX: 3, VelocityY: 4}
	if !r.DispatchGesture(id, want) || got != want {
		t.Fatalf("got %#v want %#v", got, want)
	}
	r.Release(id)
	if r.DispatchGesture(id, want) {
		t.Fatal("released gesture dispatched")
	}
}
