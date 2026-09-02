// Package ui defines the platform-independent declarative UI tree.
package ui

import "sync/atomic"

// NodeID identifies a virtual node and its corresponding native view.
type NodeID uint64

// NodeType identifies a native primitive.
type NodeType uint8

const (
	NodeView NodeType = iota + 1
	NodeText
	NodeButton
	NodeRow
	NodeColumn
	NodeSafeArea
	NodeTextInput
	NodeSwitch
	NodeProgressIndicator
	NodeImage
	NodeScrollView
)

// HandlerID identifies an event callback without passing a Go pointer to native code.
type HandlerID uint64

// AxisAlignment controls child placement on a container's main axis.
type AxisAlignment uint8

const (
	AlignStart AxisAlignment = iota
	AlignCenter
	AlignEnd
	AlignSpaceBetween
)

// AccessibilityRole describes platform-native assistive semantics.
type AccessibilityRole uint8

const (
	RoleAutomatic AccessibilityRole = iota
	RoleText
	RoleButton
	RoleHeader
	RoleImage
)

// ImageResizeMode controls how image content fits its bounds.
type ImageResizeMode uint8

const (
	ImageFit ImageResizeMode = iota
	ImageFill
	ImageCenter
)

// Props contains the compact, strongly typed properties needed by Milestone 0.
// Fields can be added without changing the shape of Node or the reconciler.
type Props struct {
	Text        string
	Width       float32
	Height      float32
	Padding     float32
	Gap         float32
	Alignment   AxisAlignment
	FontSize    float32
	Bold        bool
	OnPress     HandlerID
	OnChange    HandlerID
	OnToggle    HandlerID
	Checked     bool
	Progress    float32
	ImageSource string
	ImageMode   ImageResizeMode
	Horizontal  bool
	// Interactions is the runtime-generated comparable wire payload for gestures and animations.
	Interactions string
	AccessLabel  string
	AccessHint   string
	AccessRole   AccessibilityRole
	Focused      bool
	ScalesText   bool
}

// Node is a platform-independent native UI primitive.
type Node struct {
	ID NodeID
	// ExplicitID is set by WithID. The runtime never replaces explicit identity
	// during structural stabilization.
	ExplicitID bool
	Type       NodeType
	Props      Props
	Children   []*Node
	// Press is Go-owned behavior. It never enters Props or crosses a native boundary.
	Press func()
	// Change is Go-owned value behavior and never crosses the native boundary.
	Change func(string)
	// Toggle is Go-owned boolean behavior and never crosses the native boundary.
	Toggle            func(bool)
	Intents           IntentSet
	GestureHandlerIDs []HandlerID
}

// Component builds a virtual UI node.
type Component interface {
	Build() *Node
}

type element struct{ node *Node }

func (e *element) Build() *Node { return cloneNode(e.node) }

var nextNodeID atomic.Uint64

func newElement(kind NodeType, props Props, children ...Component) *element {
	n := &Node{ID: NodeID(nextNodeID.Add(1)), Type: kind, Props: props}
	n.Children = make([]*Node, 0, len(children))
	for _, child := range children {
		if child != nil {
			n.Children = append(n.Children, child.Build())
		}
	}
	return &element{node: n}
}

func cloneNode(n *Node) *Node {
	if n == nil {
		return nil
	}
	copy := *n
	copy.Intents = cloneIntents(n.Intents)
	copy.GestureHandlerIDs = append([]HandlerID(nil), n.GestureHandlerIDs...)
	copy.Children = make([]*Node, len(n.Children))
	for i, child := range n.Children {
		copy.Children[i] = cloneNode(child)
	}
	return &copy
}

// WithID assigns an explicit stable identity. It is primarily useful for keyed children.
func WithID(component Component, id NodeID) Component {
	n := component.Build()
	n.ID = id
	n.ExplicitID = true
	return &element{node: n}
}
