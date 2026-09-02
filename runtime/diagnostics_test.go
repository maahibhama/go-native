package runtime

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/go-native/go-native/ui"
)

func TestSnapshotTreeIsDetachedAndSerializable(t *testing.T) {
	pressed := false
	root := &ui.Node{
		ID:   1,
		Type: ui.NodeColumn,
		Children: []*ui.Node{{
			ID:    2,
			Type:  ui.NodeButton,
			Props: ui.Props{Text: "Increment", OnPress: 7},
			Press: func() { pressed = true },
		}},
	}

	snapshot := SnapshotTree(root)
	root.Children[0].Props.Text = "Changed"
	root.Children = nil

	if snapshot.Root == nil || snapshot.Root.Type != "Column" {
		t.Fatalf("unexpected root snapshot: %#v", snapshot.Root)
	}
	if got := snapshot.Root.Children[0].Props.Text; got != "Increment" {
		t.Fatalf("snapshot changed with source tree: %q", got)
	}
	encoded, err := json.Marshal(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	if pressed {
		t.Fatal("snapshotting invoked a callback")
	}
	if len(encoded) == 0 {
		t.Fatal("expected JSON output")
	}
}

func TestDiagnosticsCopiesAndClearsEntries(t *testing.T) {
	diagnostics := NewDiagnostics()
	diagnostics.Record(LogEntry{Kind: LogEventDispatched, HandlerID: 9})

	entries := diagnostics.Entries()
	if len(entries) != 1 || entries[0].Kind != LogEventDispatched {
		t.Fatalf("unexpected entries: %#v", entries)
	}
	if entries[0].Time.IsZero() {
		t.Fatal("Record did not assign a timestamp")
	}
	entries[0].Message = "mutated copy"
	if diagnostics.Entries()[0].Message != "" {
		t.Fatal("Entries exposed internal storage")
	}

	diagnostics.Clear()
	if got := diagnostics.Entries(); len(got) != 0 {
		t.Fatalf("Clear left %d entries", len(got))
	}
}

func TestDiagnosticsBoundsConcurrentHistory(t *testing.T) {
	diagnostics := NewDiagnostics()
	var group sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		group.Add(1)
		go func() {
			defer group.Done()
			for i := 0; i < diagnosticHistoryLimit; i++ {
				diagnostics.Record(LogEntry{Time: time.Unix(0, int64(i)), Kind: LogBatchApplied})
			}
		}()
	}
	group.Wait()

	if got := len(diagnostics.Entries()); got != diagnosticHistoryLimit {
		t.Fatalf("history length = %d, want %d", got, diagnosticHistoryLimit)
	}
}

func TestSnapshotTreeHandlesNil(t *testing.T) {
	if snapshot := SnapshotTree(nil); snapshot.Root != nil {
		t.Fatalf("nil tree produced root: %#v", snapshot.Root)
	}
}

func TestNodeTypeNameMapping(t *testing.T) {
	tests := []struct {
		kind ui.NodeType
		want string
	}{
		{ui.NodeView, "View"},
		{ui.NodeText, "Text"},
		{ui.NodeButton, "Button"},
		{ui.NodeRow, "Row"},
		{ui.NodeColumn, "Column"},
		{ui.NodeSafeArea, "SafeArea"},
		{ui.NodeTextInput, "TextInput"},
		{ui.NodeSwitch, "Switch"},
		{ui.NodeProgressIndicator, "ProgressIndicator"},
		{ui.NodeImage, "Image"},
		{ui.NodeScrollView, "ScrollView"},
		{ui.NodeType(255), "Unknown"},
	}

	for _, tt := range tests {
		node := &ui.Node{ID: 1, Type: tt.kind}
		snapshot := SnapshotTree(node)
		if snapshot.Root == nil || snapshot.Root.Type != tt.want {
			t.Errorf("nodeTypeName(%v) = %q, want %q", tt.kind, snapshot.Root.Type, tt.want)
		}
	}
}
