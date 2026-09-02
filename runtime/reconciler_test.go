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
	if len(got.Mutations) != 1 || got.Mutations[0].Type != MutationMove || got.Mutations[0].NodeID != 3 || got.Mutations[0].FromIndex != 1 || got.Mutations[0].Index != 0 {
		t.Fatalf("unexpected: %#v", got)
	}
}

func TestReconcileComplexOrderingUsesCurrentIndexes(t *testing.T) {
	a := node(1, ui.NodeRow, "", node(2, ui.NodeText, "a"), node(3, ui.NodeText, "b"), node(4, ui.NodeText, "c"), node(5, ui.NodeText, "d"))
	b := node(1, ui.NodeRow, "", node(5, ui.NodeText, "d"), node(3, ui.NodeText, "b"), node(6, ui.NodeText, "e"), node(2, ui.NodeText, "a"))
	batch := Reconcile(a, b)
	order := []ui.NodeID{2, 3, 4, 5}
	for _, mutation := range batch.Mutations {
		switch mutation.Type {
		case MutationRemove:
			order = append(order[:mutation.Index], order[mutation.Index+1:]...)
		case MutationInsert:
			order = insertID(order, int(mutation.Index), mutation.NodeID)
		case MutationMove:
			if order[mutation.FromIndex] != mutation.NodeID {
				t.Fatalf("move source %d contains %d, expected %d", mutation.FromIndex, order[mutation.FromIndex], mutation.NodeID)
			}
			order = moveID(order, int(mutation.FromIndex), int(mutation.Index))
		}
	}
	want := []ui.NodeID{5, 3, 6, 2}
	if !reflect.DeepEqual(order, want) {
		t.Fatalf("final order %v, want %v; mutations=%#v", order, want, batch.Mutations)
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
	in.Mutations[0].Props.AccessHint = "spoken hint"
	in.Mutations[0].Props.AccessRole = ui.RoleHeader
	in.Mutations[0].Props.Focused = true
	in.Mutations[0].Props.ScalesText = true
	in.Mutations[0].Props.OnChange = 42
	in.Mutations[0].Props.OnToggle = 43
	in.Mutations[0].Props.Checked = true
	in.Mutations[0].Props.Progress = .75
	in.Mutations[0].Props.ImageSource = "logo"
	in.Mutations[0].Props.ImageMode = ui.ImageFill
	in.Mutations[0].Props.Horizontal = true
	in.Mutations[0].Props.Interactions = string([]byte{0, 1, 0, 255})
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
