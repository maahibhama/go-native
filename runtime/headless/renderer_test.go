package headless

import (
	"testing"

	gonative "github.com/go-native/go-native/runtime"
	"github.com/go-native/go-native/ui"
)

func TestRendererAppliesCreateUpdateAndMove(t *testing.T) {
	renderer := New()
	initial := gonative.Reconcile(nil, ui.Column(ui.Text("a"), ui.Text("b")).Build())
	if err := renderer.Apply(initial); err != nil {
		t.Fatal(err)
	}
	snapshot := renderer.Snapshot()
	if snapshot.Root == nil || len(snapshot.Root.Children) != 2 || snapshot.Root.Children[0].Props.Text != "a" {
		t.Fatalf("initial snapshot = %#v", snapshot)
	}

	oldTree := ui.Column(ui.WithID(ui.Text("a"), 101), ui.WithID(ui.Text("b"), 102)).Build()
	newTree := ui.Column(ui.WithID(ui.Text("changed"), 102), ui.WithID(ui.Text("a"), 101)).Build()
	oldTree.ID = snapshot.Root.ID
	newTree.ID = snapshot.Root.ID
	// Seed a fresh renderer with explicitly keyed trees so move semantics are deterministic.
	renderer = New()
	_ = renderer.Apply(gonative.Reconcile(nil, oldTree))
	_ = renderer.Apply(gonative.Reconcile(oldTree, newTree))
	snapshot = renderer.Snapshot()
	if snapshot.Root.Children[0].ID != 102 || snapshot.Root.Children[0].Props.Text != "changed" {
		t.Fatalf("moved snapshot = %#v", snapshot)
	}
}
