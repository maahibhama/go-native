// Command gonative provides the minimal Milestone 0 developer workflow.
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const usage = `Go Native Milestone 0 CLI

Usage:
  gonative build <ios|android>
  gonative run <ios|android>
  gonative doctor
  gonative help

Environment:
  GONATIVE_SIMULATOR       iOS Simulator device name
  GONATIVE_ANDROID_SERIAL  adb device serial
  ANDROID_SDK_ROOT         Android SDK location
`

type commandRunner interface {
	Run(name string, args []string, dir string, env []string, stdout, stderr io.Writer) error
}

type processRunner struct{}

func (processRunner) Run(name string, args []string, dir string, env []string, stdout, stderr io.Writer) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

func main() {
	if err := run(os.Args[1:], processRunner{}, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "gonative:", err)
		os.Exit(1)
	}
}

func run(args []string, runner commandRunner, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		fmt.Fprint(stdout, usage)
		return nil
	}
	root, err := findProjectRoot()
	if err != nil {
		return err
	}
	switch args[0] {
	case "doctor":
		return doctor(root, stdout)
	case "build", "run":
		if len(args) != 2 {
			return fmt.Errorf("%s requires one platform\n\n%s", args[0], usage)
		}
		return platformCommand(root, args[0], args[1], runner, stdout, stderr)
	default:
		return fmt.Errorf("unknown command %q\n\n%s", args[0], usage)
	}
}

func platformCommand(root, action, platform string, runner commandRunner, stdout, stderr io.Writer) error {
	if platform != "ios" && platform != "android" {
		return fmt.Errorf("unsupported platform %q; expected ios or android", platform)
	}
	script := filepath.Join(root, "scripts", action+"-"+platform+".sh")
	if _, err := os.Stat(script); err != nil {
		return fmt.Errorf("platform script: %w", err)
	}
	env := append([]string{}, os.Environ()...)
	env = defaultEnv(env, "GOCACHE", filepath.Join(os.TempDir(), "go-native-gocache"))
	env = defaultEnv(env, "GOPATH", filepath.Join(os.TempDir(), "go-native-gopath"))
	if err := runner.Run(script, nil, root, env, stdout, stderr); err != nil {
		return fmt.Errorf("%s %s failed: %w", action, platform, err)
	}
	return nil
}

func defaultEnv(env []string, key, value string) []string {
	prefix := key + "="
	for _, item := range env {
		if strings.HasPrefix(item, prefix) {
			return env
		}
	}
	return append(env, prefix+value)
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		module, readErr := os.ReadFile(filepath.Join(dir, "go.mod"))
		if readErr == nil && strings.Contains(string(module), "module github.com/go-native/go-native") {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("not inside the go-native repository")
		}
		dir = parent
	}
}

type check struct{ label, command, path string }

func doctor(root string, out io.Writer) error {
	fmt.Fprintf(out, "Go Native doctor (%s/%s)\n", runtime.GOOS, runtime.GOARCH)
	checks := []check{
		{label: "Go", command: "go"},
		{label: "Xcode build tools", command: "xcodebuild"},
		{label: "Java compiler", command: "javac"},
		{label: "iOS build script", path: filepath.Join(root, "scripts", "build-ios.sh")},
		{label: "Android build script", path: filepath.Join(root, "scripts", "build-android.sh")},
	}
	sdk := os.Getenv("ANDROID_SDK_ROOT")
	if sdk == "" {
		sdk = os.Getenv("ANDROID_HOME")
	}
	if sdk == "" {
		home, _ := os.UserHomeDir()
		sdk = filepath.Join(home, "Library", "Android", "sdk")
	}
	checks = append(checks, check{label: "Android SDK", path: sdk})
	missing := 0
	for _, item := range checks {
		var value string
		var err error
		if item.command != "" {
			value, err = exec.LookPath(item.command)
		} else {
			value = item.path
			_, err = os.Stat(item.path)
		}
		if err != nil {
			missing++
			fmt.Fprintf(out, "  missing  %s\n", item.label)
		} else {
			fmt.Fprintf(out, "  ok       %-20s %s\n", item.label, value)
		}
	}
	if missing > 0 {
		return fmt.Errorf("%d required tool or path missing", missing)
	}
	return nil
}
