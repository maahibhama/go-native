package runtime

import (
	"github.com/go-native/go-native/ui"
	"reflect"
	"testing"
)

func node(id ui.NodeID, kind ui.NodeType, text string, children ...*ui.Node) *ui.Node {
	return &ui.Node{ID: id, Type: kind, Props: ui.Props{Text: text}, Children: children}
}
func types(b MutationBatch) []MutationType {
	v := make([]MutationType, len(b.Mutations))
	for i, m := range b.Mutations {
		v[i] = m.Type
	}
	return v
}

func TestReconcileCreation(t *testing.T) {
	b := Reconcile(nil, node(1, ui.NodeColumn, "", node(2, ui.NodeText, "hi")))
	want := []MutationType{MutationCreate, MutationCreate, MutationInsert}
	if !reflect.DeepEqual(types(b), want) {
		t.Fatalf("got %v want %v", types(b), want)
	}
}
func TestReconcileTextUpdate(t *testing.T) {
	a := node(1, ui.NodeText, "1")
	b := node(1, ui.NodeText, "2")
	got := Reconcile(a, b)
	if len(got.Mutations) != 1 || got.Mutations[0].Type != MutationUpdate || got.Mutations[0].Props.Text != "2" {
		t.Fatalf("unexpected: %#v", got)
	}
}
func TestReconcilePropertyUpdate(t *testing.T) {
	a := node(1, ui.NodeText, "x")
	b := node(1, ui.NodeText, "x")
	b.Props.FontSize = 24
	got := Reconcile(a, b)
	if len(got.Mutations) != 1 || got.Mutations[0].Props.FontSize != 24 {
		t.Fatalf("unexpected: %#v", got)
	}
}
func TestReconcileChildInsertion(t *testing.T) {
	a := node(1, ui.NodeColumn, "", node(2, ui.NodeText, "a"))
	b := node(1, ui.NodeColumn, "", node(2, ui.NodeText, "a"), node(3, ui.NodeText, "b"))
	if got := types(Reconcile(a, b)); !reflect.DeepEqual(got, []MutationType{MutationCreate, MutationInsert}) {
		t.Fatalf("got %v", got)
	}
}
func TestReconcileChildRemoval(t *testing.T) {
	a := node(1, ui.NodeColumn, "", node(2, ui.NodeText, "a"), node(3, ui.NodeText, "b"))
	b := node(1, ui.NodeColumn, "", node(2, ui.NodeText, "a"))
	if got := types(Reconcile(a, b)); !reflect.DeepEqual(got, []MutationType{MutationRemove, MutationDelete}) {
		t.Fatalf("got %v", got)
	}
}
func TestReconcileChildOrdering(t *testing.T) {
	a := node(1, ui.NodeRow, "", node(2, ui.NodeText, "a"), node(3, ui.NodeText, "b"))
	b := node(1, ui.NodeRow, "", node(3, ui.NodeText, "b"), node(2, ui.NodeText, "a"))
	got := Reconcile(a, b)
	if len(got.Mutations) != 2 || got.Mutations[0].Type != MutationMove || got.Mutations[1].Type != MutationMove {
		t.Fatalf("unexpected: %#v", got)
	}
}
func TestReconcileNestedUpdate(t *testing.T) {
	a := node(1, ui.NodeColumn, "", node(2, ui.NodeRow, "", node(3, ui.NodeText, "old")))
	b := node(1, ui.NodeColumn, "", node(2, ui.NodeRow, "", node(3, ui.NodeText, "new")))
	got := Reconcile(a, b)
	if len(got.Mutations) != 1 || got.Mutations[0].NodeID != 3 {
		t.Fatalf("unexpected: %#v", got)
	}
}

func TestMutationBatchBinaryRoundTrip(t *testing.T) {
	in := Reconcile(nil, node(1, ui.NodeText, "hello"))
	data, err := in.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	out, err := UnmarshalMutationBatch(data)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(in, out) {
		t.Fatalf("got %#v want %#v", out, in)
	}
}
