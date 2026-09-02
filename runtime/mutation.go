// Package runtime reconciles virtual trees and drives platform renderers.
package runtime

import "github.com/go-native/go-native/ui"

import "time"

// MutationType identifies an operation applied by a native renderer.
type MutationType uint8

const (
	MutationCreate MutationType = iota + 1
	MutationDelete
	MutationUpdate
	MutationInsert
	MutationRemove
	MutationMove
)

// Mutation is one typed change to the native view tree.
type Mutation struct {
	Type      MutationType
	NodeID    ui.NodeID
	ParentID  ui.NodeID
	NodeType  ui.NodeType
	Index     int32
	FromIndex int32
	Props     ui.Props
}

// MutationBatch crosses the native boundary as one serialized payload.
type MutationBatch struct {
	Sequence  uint64
	Mutations []Mutation
}

// TimingSample describes one batch acknowledged by a native renderer.
type TimingSample struct {
	Sequence      uint64
	MutationCount int
	NativeApply   time.Duration
	BridgeToApply time.Duration
	EventToApply  time.Duration
}

// Renderer accepts complete batches. Implementations must marshal application to the UI thread.
type Renderer interface{ Apply(MutationBatch) error }
