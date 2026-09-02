package runtime

import (
	"strconv"
	"testing"

	"github.com/go-native/go-native/ui"
)

func benchmarkTree(size int) *ui.Node {
	children := make([]*ui.Node, size)
	for i := range children {
		children[i] = node(ui.NodeID(i+2), ui.NodeText, strconv.Itoa(i))
	}
	return node(1, ui.NodeColumn, "", children...)
}
func benchmarkDeclarativeTree(size int) ui.Component {
	children := make([]ui.Component, size)
	for i := range children {
		children[i] = ui.Text(strconv.Itoa(i))
	}
	return ui.Column(children...)
}

func BenchmarkDeclarativeTreeBuild1000(b *testing.B) {
	for range b.N {
		benchmarkDeclarativeTree(1000).Build()
	}
}

func BenchmarkInitialReconcile1000(b *testing.B) {
	tree := benchmarkTree(1000)
	b.ResetTimer()
	for range b.N {
		Reconcile(nil, tree)
	}
}

func BenchmarkReconcileOnePropertyUpdate1000(b *testing.B) {
	old := benchmarkTree(1000)
	next := benchmarkTree(1000)
	next.Children[500].Props.Text = "changed"
	b.ResetTimer()
	for range b.N {
		Reconcile(old, next)
	}
}

func BenchmarkReconcile100PropertyUpdates1000(b *testing.B) {
	old := benchmarkTree(1000)
	next := benchmarkTree(1000)
	for i := 0; i < 100; i++ {
		next.Children[i*10].Props.Text = "changed"
	}
	b.ResetTimer()
	for range b.N {
		Reconcile(old, next)
	}
}

func BenchmarkReconcileNoChanges1000(b *testing.B) {
	old := benchmarkTree(1000)
	next := benchmarkTree(1000)
	b.ResetTimer()
	for range b.N {
		Reconcile(old, next)
	}
}

func BenchmarkSerialize1000Creates(b *testing.B) {
	batch := Reconcile(nil, benchmarkTree(1000))
	b.ResetTimer()
	for range b.N {
		_, _ = batch.MarshalBinary()
	}
}

func BenchmarkSerializeOneUpdate(b *testing.B) {
	old := benchmarkTree(1)
	next := benchmarkTree(1)
	next.Children[0].Props.Text = "changed"
	batch := Reconcile(old, next)
	b.ResetTimer()
	for range b.N {
		_, _ = batch.MarshalBinary()
	}
}

func BenchmarkSerialize100Updates(b *testing.B) {
	old := benchmarkTree(100)
	next := benchmarkTree(100)
	for i := range next.Children {
		next.Children[i].Props.Text = "changed"
	}
	batch := Reconcile(old, next)
	b.ResetTimer()
	for range b.N {
		_, _ = batch.MarshalBinary()
	}
}

func BenchmarkEventRegistryDispatch(b *testing.B) {
	registry := NewEventRegistry()
	id := registry.Register(func() {})
	b.ResetTimer()
	for range b.N {
		registry.Dispatch(id)
	}
}
