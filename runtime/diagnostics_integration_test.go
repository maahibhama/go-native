package runtime

import (
	"testing"

	"github.com/go-native/go-native/ui"
)

func TestRuntimeDiagnosticsAndTreeSnapshot(t *testing.T) {
	renderer := &recordingRenderer{}
	r := New(func() ui.Component { return ui.Column(ui.Text("visible")) }, renderer)
	if err := r.Start(); err != nil {
		t.Fatal(err)
	}
	entries := r.LogEntries()
	if len(entries) != 2 || entries[0].Kind != LogRuntimeStarted || entries[1].Kind != LogBatchApplied {
		t.Fatalf("unexpected startup log: %#v", entries)
	}
	snapshot := r.TreeSnapshot()
	if snapshot.Root == nil || snapshot.Root.Type != "Column" || len(snapshot.Root.Children) != 1 || snapshot.Root.Children[0].Props.Text != "visible" {
		t.Fatalf("unexpected tree snapshot: %#v", snapshot)
	}
	if r.Dispatch(ui.HandlerID(999)) {
		t.Fatal("missing event dispatched")
	}
	entries = r.LogEntries()
	if entries[len(entries)-1].Kind != LogEventMissing {
		t.Fatalf("missing event was not logged: %#v", entries)
	}
	r.ClearLogEntries()
	if len(r.LogEntries()) != 0 {
		t.Fatal("diagnostics did not clear")
	}
	r.Stop()
}
