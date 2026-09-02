package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type recordedRun struct {
	name, dir string
	args, env []string
}
type fakeRunner struct {
	calls []recordedRun
	err   error
}

func (f *fakeRunner) Run(name string, args []string, dir string, env []string, stdout, stderr io.Writer) error {
	f.calls = append(f.calls, recordedRun{name: name, args: args, dir: dir, env: env})
	return f.err
}

func TestHelp(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"help"}, &fakeRunner{}, &out, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "gonative build") {
		t.Fatalf("unexpected help: %s", out.String())
	}
}

func TestBuildIOSDispatchesScript(t *testing.T) {
	root, err := findProjectRoot()
	if err != nil {
		t.Fatal(err)
	}
	runner := &fakeRunner{}
	var out bytes.Buffer
	if err = run([]string{"build", "ios"}, runner, &out, &out); err != nil {
		t.Fatal(err)
	}
	if len(runner.calls) != 1 {
		t.Fatalf("calls=%d", len(runner.calls))
	}
	want := filepath.Join(root, "scripts", "build-ios.sh")
	if runner.calls[0].name != want {
		t.Fatalf("script=%s want=%s", runner.calls[0].name, want)
	}
}

func TestUnsupportedPlatform(t *testing.T) {
	var out bytes.Buffer
	err := run([]string{"build", "web"}, &fakeRunner{}, &out, &out)
	if err == nil || !strings.Contains(err.Error(), "unsupported platform") {
		t.Fatalf("error=%v", err)
	}
}

func TestDefaultEnvPreservesValue(t *testing.T) {
	env := defaultEnv([]string{"GOCACHE=/custom"}, "GOCACHE", "/tmp/default")
	if len(env) != 1 || env[0] != "GOCACHE=/custom" {
		t.Fatalf("env=%v", env)
	}
}

func TestFindProjectRootFromChild(t *testing.T) {
	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root, err := findProjectRoot()
	if err != nil {
		t.Fatal(err)
	}
	if err = os.Chdir(filepath.Join(root, "ui")); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(original)
	got, err := findProjectRoot()
	if err != nil {
		t.Fatal(err)
	}
	if got != root {
		t.Fatalf("got=%s want=%s", got, root)
	}
}
