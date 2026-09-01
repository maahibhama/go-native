package runtime

import (
	"github.com/go-native/go-native/ui"
	"sync"
	"testing"
	"time"
)

type recordingRenderer struct {
	mu      sync.Mutex
	batches []MutationBatch
}

func (r *recordingRenderer) Apply(b MutationBatch) error {
	r.mu.Lock()
	r.batches = append(r.batches, b)
	r.mu.Unlock()
	return nil
}

func TestRuntimeStateChangeOnlyUpdatesText(t *testing.T) {
	state := ui.NewState(0)
	renderer := &recordingRenderer{}
	app := func() ui.Component {
		return ui.Column(ui.Text(string(rune('0'+state.Get()))), ui.Button("Increment", func() { state.Set(state.Get() + 1) }))
	}
	rt := New(app, renderer)
	if err := rt.Start(); err != nil {
		t.Fatal(err)
	}
	renderer.mu.Lock()
	handler := renderer.batches[0].Mutations[len(renderer.batches[0].Mutations)-2].Props.OnPress
	renderer.mu.Unlock()
	if handler == 0 {
		t.Fatal("button handler not bound")
	}
	rt.Dispatch(handler)
	deadline := time.Now().Add(time.Second)
	for {
		renderer.mu.Lock()
		n := len(renderer.batches)
		renderer.mu.Unlock()
		if n >= 2 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("render not scheduled")
		}
		time.Sleep(time.Millisecond)
	}
	renderer.mu.Lock()
	defer renderer.mu.Unlock()
	got := renderer.batches[1].Mutations
	if len(got) != 1 || got[0].Type != MutationUpdate || got[0].NodeType != ui.NodeText || got[0].Props.Text != "1" {
		t.Fatalf("expected one text update, got %#v", got)
	}
}
