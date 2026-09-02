package runtime

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/go-native/go-native/ui"
)

// LogKind identifies a portable runtime diagnostic event.
type LogKind string

const (
	LogRuntimeStarted  LogKind = "runtime.started"
	LogRuntimeStopped  LogKind = "runtime.stopped"
	LogEventDispatched LogKind = "event.dispatched"
	LogEventMissing    LogKind = "event.missing"
	LogRenderSkipped   LogKind = "render.skipped"
	LogBatchApplied    LogKind = "render.batch_applied"
	LogRenderFailed    LogKind = "render.failed"
)

// LogEntry is a structured, platform-independent runtime event. Sequence and
// MutationCount are populated for renderer batches; HandlerID is populated for
// event dispatches.
type LogEntry struct {
	Time          time.Time    `json:"time"`
	Kind          LogKind      `json:"kind"`
	Sequence      uint64       `json:"sequence,omitempty"`
	MutationCount int          `json:"mutationCount,omitempty"`
	HandlerID     ui.HandlerID `json:"handlerId,omitempty"`
	Message       string       `json:"message,omitempty"`
}

// TreeSnapshot is an immutable representation of the last successfully
// rendered virtual tree. It is safe to encode or inspect on another goroutine.
type TreeSnapshot struct {
	CapturedAt time.Time         `json:"capturedAt"`
	Root       *TreeNodeSnapshot `json:"root,omitempty"`
}

// TreeNodeSnapshot is the serializable form of a virtual UI node. Go callback
// functions and internal identity metadata are intentionally omitted.
type TreeNodeSnapshot struct {
	ID       ui.NodeID          `json:"id"`
	Type     string             `json:"type"`
	Props    ui.Props           `json:"props"`
	Children []TreeNodeSnapshot `json:"children,omitempty"`
}

// MarshalJSON returns a JSON document suitable for a debugger or inspector.
func (s TreeSnapshot) MarshalJSON() ([]byte, error) {
	type snapshotAlias TreeSnapshot
	return json.Marshal(snapshotAlias(s))
}

// Diagnostics stores a bounded event history and produces detached UI-tree
// snapshots. A Runtime can own one without involving either native renderer.
type Diagnostics struct {
	mu      sync.Mutex
	entries []LogEntry
}

const diagnosticHistoryLimit = 1024

// NewDiagnostics creates an empty portable diagnostics collector.
func NewDiagnostics() *Diagnostics { return &Diagnostics{} }

// Record appends a structured event. Zero timestamps are filled at ingestion.
func (d *Diagnostics) Record(entry LogEntry) {
	if entry.Time.IsZero() {
		entry.Time = time.Now()
	}
	d.mu.Lock()
	d.entries = append(d.entries, entry)
	if len(d.entries) > diagnosticHistoryLimit {
		copy(d.entries, d.entries[len(d.entries)-diagnosticHistoryLimit:])
		d.entries = d.entries[:diagnosticHistoryLimit]
	}
	d.mu.Unlock()
}

// Entries returns a chronological copy of the bounded event history.
func (d *Diagnostics) Entries() []LogEntry {
	d.mu.Lock()
	defer d.mu.Unlock()
	return append([]LogEntry(nil), d.entries...)
}

// Clear removes all accumulated entries.
func (d *Diagnostics) Clear() {
	d.mu.Lock()
	d.entries = nil
	d.mu.Unlock()
}

// SnapshotTree returns a detached, serializable copy of root. The caller must
// provide synchronization if the source tree can be mutated concurrently.
func SnapshotTree(root *ui.Node) TreeSnapshot {
	return TreeSnapshot{CapturedAt: time.Now(), Root: snapshotNode(root)}
}

func snapshotNode(n *ui.Node) *TreeNodeSnapshot {
	if n == nil {
		return nil
	}
	snapshot := &TreeNodeSnapshot{ID: n.ID, Type: nodeTypeName(n.Type), Props: n.Props}
	if len(n.Children) > 0 {
		snapshot.Children = make([]TreeNodeSnapshot, len(n.Children))
		for i, child := range n.Children {
			snapshot.Children[i] = *snapshotNode(child)
		}
	}
	return snapshot
}

func nodeTypeName(kind ui.NodeType) string {
	switch kind {
	case ui.NodeView:
		return "View"
	case ui.NodeText:
		return "Text"
	case ui.NodeButton:
		return "Button"
	case ui.NodeRow:
		return "Row"
	case ui.NodeColumn:
		return "Column"
	case ui.NodeSafeArea:
		return "SafeArea"
	case ui.NodeTextInput:
		return "TextInput"
	case ui.NodeSwitch:
		return "Switch"
	case ui.NodeProgressIndicator:
		return "ProgressIndicator"
	case ui.NodeImage:
		return "Image"
	case ui.NodeScrollView:
		return "ScrollView"
	default:
		return "Unknown"
	}
}
