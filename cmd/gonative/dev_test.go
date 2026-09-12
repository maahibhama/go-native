package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestScanSourcesFiltersGeneratedDirectories(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "app.go"), "package app")
	writeTestFile(t, filepath.Join(root, "assets", "icon.txt"), "one")
	writeTestFile(t, filepath.Join(root, "build", "generated.go"), "ignored")
	writeTestFile(t, filepath.Join(root, ".build", "generated.go"), "ignored")
	snapshot, err := scanSources(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot) != 2 {
		t.Fatalf("snapshot = %v", snapshot)
	}
}

func TestDiffSnapshotsReportsCreateChangeAndDelete(t *testing.T) {
	root := t.TempDir()
	one := filepath.Join(root, "one.go")
	two := filepath.Join(root, "two.go")
	writeTestFile(t, one, "one")
	before, _ := scanSources(root)
	writeTestFile(t, one, "changed")
	writeTestFile(t, two, "two")
	after, _ := scanSources(root)
	changes := diffSnapshots(before, after)
	if len(changes) != 2 || changes[0] != "one.go" || changes[1] != "two.go" {
		t.Fatalf("changes = %v", changes)
	}
	if err := os.Remove(one); err != nil {
		t.Fatal(err)
	}
	deleted, _ := scanSources(root)
	changes = diffSnapshots(after, deleted)
	if len(changes) != 1 || changes[0] != "one.go" {
		t.Fatalf("deleted changes = %v", changes)
	}
}

func TestWatchSourcesStopsAndReportsChange(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "app.go")
	writeTestFile(t, path, "one")
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- watchSources(ctx, root, time.Millisecond, func(changes []string) { cancel() })
	}()
	time.Sleep(5 * time.Millisecond)
	writeTestFile(t, path, "two")
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("watcher did not stop")
	}
}

func TestLoadDevSessionReusesAndResets(t *testing.T) {
	root := t.TempDir()
	first, err := loadDevSession(root, false)
	if err != nil {
		t.Fatal(err)
	}
	second, err := loadDevSession(root, false)
	if err != nil || second != first {
		t.Fatalf("reused session = %q, %v; want %q", second, err, first)
	}
	reset, err := loadDevSession(root, true)
	if err != nil || reset == first {
		t.Fatalf("reset session = %q, %v; previous %q", reset, err, first)
	}
}

func TestParseDevCommand(t *testing.T) {
	tests := map[string]string{"i": "ios", "IOS": "ios", "a": "android", "r": "reload", "d": "doctor", "q": "quit", "unknown": "help"}
	for input, want := range tests {
		if got := parseDevCommand(input); got != want {
			t.Errorf("parseDevCommand(%q) = %q, want %q", input, got, want)
		}
	}
}

func writeTestFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
