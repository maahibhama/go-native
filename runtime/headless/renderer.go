// Package headless provides a deterministic in-memory native host for tests.
package headless

import (
	"sync"

	gonative "github.com/go-native/go-native/runtime"
	"github.com/go-native/go-native/ui"
)

type Renderer struct {
	mu      sync.Mutex
	root    *ui.Node
	nodes   map[ui.NodeID]*ui.Node
	batches []gonative.MutationBatch
}

func New() *Renderer { return &Renderer{nodes: make(map[ui.NodeID]*ui.Node)} }

func (r *Renderer) Apply(batch gonative.MutationBatch) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.nodes == nil {
		r.nodes = make(map[ui.NodeID]*ui.Node)
	}
	for _, mutation := range batch.Mutations {
		r.apply(mutation)
	}
	r.batches = append(r.batches, cloneBatch(batch))
	return nil
}

func (r *Renderer) apply(m gonative.Mutation) {
	switch m.Type {
	case gonative.MutationCreate:
		n := &ui.Node{ID: m.NodeID, Type: m.NodeType, Props: m.Props}
		r.nodes[m.NodeID] = n
		if r.root == nil {
			r.root = n
		}
	case gonative.MutationUpdate:
		if n := r.nodes[m.NodeID]; n != nil {
			n.Props = m.Props
		}
	case gonative.MutationInsert:
		parent, child := r.nodes[m.ParentID], r.nodes[m.NodeID]
		if parent == nil || child == nil {
			return
		}
		index := clampIndex(int(m.Index), len(parent.Children))
		parent.Children = append(parent.Children, nil)
		copy(parent.Children[index+1:], parent.Children[index:])
		parent.Children[index] = child
	case gonative.MutationRemove:
		if parent := r.nodes[m.ParentID]; parent != nil {
			parent.Children = remove(parent.Children, m.NodeID)
		}
	case gonative.MutationMove:
		if parent := r.nodes[m.ParentID]; parent != nil {
			child := r.nodes[m.NodeID]
			parent.Children = remove(parent.Children, m.NodeID)
			index := clampIndex(int(m.Index), len(parent.Children))
			parent.Children = append(parent.Children, nil)
			copy(parent.Children[index+1:], parent.Children[index:])
			parent.Children[index] = child
		}
	case gonative.MutationDelete:
		if r.root != nil && r.root.ID == m.NodeID {
			r.root = nil
		}
		delete(r.nodes, m.NodeID)
	}
}

func (r *Renderer) Snapshot() gonative.TreeSnapshot {
	r.mu.Lock()
	defer r.mu.Unlock()
	return gonative.SnapshotTree(r.root)
}

func (r *Renderer) Batches() []gonative.MutationBatch {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]gonative.MutationBatch, len(r.batches))
	for i, batch := range r.batches {
		out[i] = cloneBatch(batch)
	}
	return out
}

func cloneBatch(batch gonative.MutationBatch) gonative.MutationBatch {
	batch.Mutations = append([]gonative.Mutation(nil), batch.Mutations...)
	return batch
}
func clampIndex(index, length int) int {
	if index < 0 {
		return 0
	}
	if index > length {
		return length
	}
	return index
}
func remove(children []*ui.Node, id ui.NodeID) []*ui.Node {
	for i, child := range children {
		if child != nil && child.ID == id {
			return append(children[:i], children[i+1:]...)
		}
	}
	return children
}
