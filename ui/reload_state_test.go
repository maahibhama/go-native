package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestReloadStateRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	ConfigureReloadState(ReloadStateOptions{Path: path, SessionID: "test"})
	value, persist := restoreReloadValue("screen:count", 1)
	if value != 1 || persist == nil {
		t.Fatalf("initial = %d, persist nil = %v", value, persist == nil)
	}
	persist(7)
	waitForReloadFile(t, path)
	ConfigureReloadState(ReloadStateOptions{Path: path, SessionID: "test"})
	value, _ = restoreReloadValue("screen:count", 1)
	if value != 7 {
		t.Fatalf("restored = %d, want 7", value)
	}
	ConfigureReloadState(ReloadStateOptions{})
}

func TestReloadStateRejectsSessionAndCorruptValues(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	ConfigureReloadState(ReloadStateOptions{Path: path, SessionID: "old"})
	_, persist := restoreReloadValue("screen:count", 0)
	persist(4)
	waitForReloadFile(t, path)
	var warnings []string
	ConfigureReloadState(ReloadStateOptions{Path: path, SessionID: "new", OnWarning: func(message string) { warnings = append(warnings, message) }})
	value, _ := restoreReloadValue("screen:count", 9)
	if value != 9 || len(warnings) == 0 || !strings.Contains(warnings[0], "incompatible") {
		t.Fatalf("value=%d warnings=%v", value, warnings)
	}
	ConfigureReloadState(ReloadStateOptions{})
}

func waitForReloadFile(t *testing.T, path string) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for {
		if _, err := os.Stat(path); err == nil {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("reload state was not written: %s", path)
		}
		time.Sleep(time.Millisecond)
	}
}
