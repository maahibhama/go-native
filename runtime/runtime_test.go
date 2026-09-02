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

func TestStabilizeIDsPreservesExplicitIdentityAcrossReorder(t *testing.T) {
	oldTree := ui.Column(
		ui.WithID(ui.Text("a"), 100),
		ui.WithID(ui.Text("b"), 200),
	).Build()
	newTree := ui.Column(
		ui.WithID(ui.Text("b"), 200),
		ui.WithID(ui.Text("a"), 100),
	).Build()
	stabilizeIDs(oldTree, newTree)
	if newTree.ID != oldTree.ID {
		t.Fatal("unkeyed root identity was not stabilized")
	}
	if newTree.Children[0].ID != 200 || newTree.Children[1].ID != 100 {
		t.Fatalf("explicit IDs were replaced: %d %d", newTree.Children[0].ID, newTree.Children[1].ID)
	}
	batch := Reconcile(oldTree, newTree)
	if len(batch.Mutations) != 1 || batch.Mutations[0].Type != MutationMove {
		t.Fatalf("expected one native move, got %#v", batch.Mutations)
	}
}

func TestKeyedButtonReorderPreservesHandlerIdentity(t *testing.T) {
	r := &Runtime{events: NewEventRegistry()}
	firstCalls, secondCalls := 0, 0
	oldTree := ui.Column(
		ui.WithID(ui.Button("first", func() { firstCalls++ }), 100),
		ui.WithID(ui.Button("second", func() { secondCalls++ }), 200),
	).Build()
	r.bindHandlers(nil, oldTree)
	firstHandler := oldTree.Children[0].Props.OnPress
	secondHandler := oldTree.Children[1].Props.OnPress

	nextTree := ui.Column(
		ui.WithID(ui.Button("second", func() { secondCalls += 10 }), 200),
		ui.WithID(ui.Button("first", func() { firstCalls += 10 }), 100),
	).Build()
	r.bindHandlers(oldTree, nextTree)
	if nextTree.Children[0].Props.OnPress != secondHandler || nextTree.Children[1].Props.OnPress != firstHandler {
		t.Fatalf("handlers followed position rather than identity: %#v", nextTree.Children)
	}
	r.events.Dispatch(nextTree.Children[0].Props.OnPress)
	r.events.Dispatch(nextTree.Children[1].Props.OnPress)
	if firstCalls != 10 || secondCalls != 10 {
		t.Fatalf("wrong callbacks dispatched: first=%d second=%d", firstCalls, secondCalls)
	}
}
