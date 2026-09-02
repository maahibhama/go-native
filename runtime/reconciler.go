package runtime

import "github.com/go-native/go-native/ui"

// Reconcile computes native mutations from oldTree to newTree.
func Reconcile(oldTree, newTree *ui.Node) MutationBatch {
	var out []Mutation
	reconcileNode(&out, 0, 0, oldTree, newTree)
	return MutationBatch{Mutations: out}
}

func reconcileNode(out *[]Mutation, parent ui.NodeID, index int, oldNode, newNode *ui.Node) {
	switch {
	case oldNode == nil && newNode != nil:
		createSubtree(out, parent, index, newNode)
	case oldNode != nil && newNode == nil:
		deleteSubtree(out, parent, index, oldNode)
	case oldNode.ID != newNode.ID || oldNode.Type != newNode.Type:
		deleteSubtree(out, parent, index, oldNode)
		createSubtree(out, parent, index, newNode)
	default:
		if oldNode.Props != newNode.Props {
			*out = append(*out, Mutation{Type: MutationUpdate, NodeID: newNode.ID, NodeType: newNode.Type, Props: newNode.Props})
		}
		reconcileChildren(out, newNode.ID, oldNode.Children, newNode.Children)
	}
}

func createSubtree(out *[]Mutation, parent ui.NodeID, index int, node *ui.Node) {
	*out = append(*out, Mutation{Type: MutationCreate, NodeID: node.ID, NodeType: node.Type, Props: node.Props})
	if parent != 0 {
		*out = append(*out, Mutation{Type: MutationInsert, NodeID: node.ID, ParentID: parent, Index: int32(index)})
	}
	for i, child := range node.Children {
		createSubtree(out, node.ID, i, child)
	}
}

func deleteSubtree(out *[]Mutation, parent ui.NodeID, index int, node *ui.Node) {
	for i := len(node.Children) - 1; i >= 0; i-- {
		deleteSubtree(out, node.ID, i, node.Children[i])
	}
	if parent != 0 {
		*out = append(*out, Mutation{Type: MutationRemove, NodeID: node.ID, ParentID: parent, Index: int32(index)})
	}
	*out = append(*out, Mutation{Type: MutationDelete, NodeID: node.ID})
}

func reconcileChildren(out *[]Mutation, parent ui.NodeID, oldChildren, newChildren []*ui.Node) {
	oldByID := make(map[ui.NodeID]*ui.Node, len(oldChildren))
	working := make([]ui.NodeID, len(oldChildren))
	for i, child := range oldChildren {
		oldByID[child.ID] = child
		working[i] = child.ID
	}
	newIDs := make(map[ui.NodeID]bool, len(newChildren))
	for _, child := range newChildren {
		newIDs[child.ID] = true
	}
	for i := len(oldChildren) - 1; i >= 0; i-- {
		if !newIDs[oldChildren[i].ID] {
			deleteSubtree(out, parent, i, oldChildren[i])
			working = append(working[:i], working[i+1:]...)
		}
	}
	for i, child := range newChildren {
		oldChild, exists := oldByID[child.ID]
		if !exists {
			createSubtree(out, parent, i, child)
			working = insertID(working, i, child.ID)
			continue
		}
		current := indexOf(working, child.ID)
		if current != i {
			*out = append(*out, Mutation{Type: MutationMove, NodeID: child.ID, ParentID: parent, FromIndex: int32(current), Index: int32(i)})
			working = moveID(working, current, i)
		}
		reconcileNode(out, parent, i, oldChild, child)
	}
}

func indexOf(ids []ui.NodeID, id ui.NodeID) int {
	for i, candidate := range ids {
		if candidate == id {
			return i
		}
	}
	return -1
}

func insertID(ids []ui.NodeID, index int, id ui.NodeID) []ui.NodeID {
	ids = append(ids, 0)
	copy(ids[index+1:], ids[index:])
	ids[index] = id
	return ids
}

func moveID(ids []ui.NodeID, from, to int) []ui.NodeID {
	id := ids[from]
	ids = append(ids[:from], ids[from+1:]...)
	return insertID(ids, to, id)
}
