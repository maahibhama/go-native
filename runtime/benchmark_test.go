package runtime

import (
	"github.com/go-native/go-native/ui"
	"strconv"
	"testing"
)

func benchmarkTree(size int) *ui.Node {
	children := make([]*ui.Node, size)
	for i := range children {
		children[i] = node(ui.NodeID(i+2), ui.NodeText, strconv.Itoa(i))
	}
	return node(1, ui.NodeColumn, "", children...)
}
func BenchmarkInitialTree1000(b *testing.B) {
	tree := benchmarkTree(1000)
	b.ResetTimer()
	for range b.N {
		Reconcile(nil, tree)
	}
}
func BenchmarkOnePropertyUpdate1000(b *testing.B) {
	old := benchmarkTree(1000)
	next := benchmarkTree(1000)
	next.Children[500].Props.Text = "changed"
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
