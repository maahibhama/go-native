package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	gnruntime "github.com/go-native/go-native/runtime"
	"github.com/go-native/go-native/runtime/devtransport"
)

const devPollInterval = 200 * time.Millisecond

type sourceSnapshot map[string][32]byte

func devStartCommand(root string, runner commandRunner, input io.Reader, stdout, stderr io.Writer) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	commands := make(chan string)
	go readDevCommands(ctx, input, commands)
	changes := make(chan []string, 1)
	go func() {
		_ = watchSources(ctx, root, devPollInterval, func(paths []string) {
			select {
			case changes <- paths:
			default:
			}
		})
	}()

	fmt.Fprintln(stdout, "Go Native development server")
	printDevMenu(stdout)
	activePlatform := ""
	builtPlatforms := make(map[string]bool)
	generations := make(map[string]uint64)
	for {
		select {
		case <-ctx.Done():
			return nil
		case command, ok := <-commands:
			if !ok {
				return nil
			}
			switch parseDevCommand(command) {
			case "ios":
				if runDevPlatform(root, "ios", builtPlatforms["ios"], runner, stdout, stderr) {
					activePlatform = "ios"
					builtPlatforms["ios"] = true
					generations["ios"]++
				}
			case "android":
				if runDevPlatform(root, "android", builtPlatforms["android"], runner, stdout, stderr) {
					activePlatform = "android"
					builtPlatforms["android"] = true
					generations["android"]++
				}
			case "reload":
				if activePlatform == "" {
					fmt.Fprintln(stdout, "No active platform. Press i for iOS or a for Android first.")
				} else {
					fmt.Fprintln(stdout, "Virtual Reload: no native-shell connection; using rebuild fallback")
					if runDevPlatform(root, activePlatform, builtPlatforms[activePlatform], runner, stdout, stderr) {
						generations[activePlatform]++
					}
				}
			case "rebuild":
				if activePlatform == "" {
					fmt.Fprintln(stdout, "No active platform. Press i for iOS or a for Android first.")
				} else if runDevPlatform(root, activePlatform, false, runner, stdout, stderr) {
					builtPlatforms[activePlatform] = true
					generations[activePlatform]++
				}
			case "clear":
				if _, err := loadDevSession(root, true); err != nil {
					fmt.Fprintln(stderr, "Clear reload state:", err)
				} else if activePlatform == "" {
					fmt.Fprintln(stdout, "Reload state cleared. Select a platform to start a fresh session.")
				} else if runDevPlatform(root, activePlatform, builtPlatforms[activePlatform], runner, stdout, stderr) {
					generations[activePlatform]++
				}
			case "status":
				printDevStatus(stdout, activePlatform, builtPlatforms, generations)
			case "inspector":
				fmt.Fprintln(stdout, "Inspector: waiting for the virtual-runtime worker connection")
			case "doctor":
				if err := doctor(root, stdout); err != nil {
					fmt.Fprintln(stderr, "doctor:", err)
				}
			case "quit":
				return nil
			case "help":
				printDevMenu(stdout)
			}
		case paths := <-changes:
			if activePlatform != "" {
				fmt.Fprintf(stdout, "Fast Reload: %s changed\n", strings.Join(paths, ", "))
				runDevPlatform(root, activePlatform, builtPlatforms[activePlatform], runner, stdout, stderr)
			}
		}
	}
}

func readDevCommands(ctx context.Context, input io.Reader, commands chan<- string) {
	defer close(commands)
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		select {
		case commands <- scanner.Text():
		case <-ctx.Done():
			return
		}
	}
}

func parseDevCommand(input string) string {
	trimmed := strings.TrimSpace(input)
	if trimmed == "R" {
		return "rebuild"
	}
	switch strings.ToLower(trimmed) {
	case "i", "ios":
		return "ios"
	case "a", "android":
		return "android"
	case "r", "reload":
		return "reload"
	case "b", "build", "rebuild":
		return "rebuild"
	case "c", "clear":
		return "clear"
	case "d", "doctor":
		return "doctor"
	case "p", "status":
		return "status"
	case "o", "open", "inspector":
		return "inspector"
	case "q", "quit", "exit":
		return "quit"
	default:
		return "help"
	}
}

func printDevMenu(output io.Writer) {
	fmt.Fprintln(output, "  i  run iOS Simulator")
	fmt.Fprintln(output, "  a  run Android emulator")
	fmt.Fprintln(output, "  r  virtual reload connected targets (rebuild fallback until connected)")
	fmt.Fprintln(output, "  R/b rebuild and reinstall the active native shell")
	fmt.Fprintln(output, "  c  clear reload state and reload")
	fmt.Fprintln(output, "  d  run doctor")
	fmt.Fprintln(output, "  p  print connection and protocol status")
	fmt.Fprintln(output, "  o  open inspector when a worker is connected")
	fmt.Fprintln(output, "  q  quit")
}

func printDevStatus(output io.Writer, active string, built map[string]bool, generations map[string]uint64) {
	if active == "" {
		active = "none"
	}
	fmt.Fprintf(output, "Development status: active=%s transport=v%d mutation=v%d\n", active, devtransport.Version, gnruntime.ProtocolVersion())
	for _, platform := range []string{"ios", "android"} {
		state := "not built"
		if built[platform] {
			state = "installed; embedded fallback"
		}
		fmt.Fprintf(output, "  %s: %s, generation=%d\n", platform, state, generations[platform])
	}
}

func runDevPlatform(root, platform string, skipFramework bool, runner commandRunner, stdout, stderr io.Writer) bool {
	session, err := loadDevSession(root, false)
	if err != nil {
		fmt.Fprintln(stderr, "Fast Reload:", err)
		return false
	}
	previousSession, hadSession := os.LookupEnv("GONATIVE_RELOAD_SESSION")
	previousSkip, hadSkip := os.LookupEnv("GONATIVE_SKIP_FRAMEWORK_BUILD")
	defer restoreEnv("GONATIVE_RELOAD_SESSION", previousSession, hadSession)
	defer restoreEnv("GONATIVE_SKIP_FRAMEWORK_BUILD", previousSkip, hadSkip)
	_ = os.Setenv("GONATIVE_RELOAD_SESSION", session)
	if skipFramework {
		_ = os.Setenv("GONATIVE_SKIP_FRAMEWORK_BUILD", "1")
	} else {
		_ = os.Unsetenv("GONATIVE_SKIP_FRAMEWORK_BUILD")
	}
	started := time.Now()
	if err := platformCommand(root, "run", platform, runner, stdout, stderr); err != nil {
		fmt.Fprintf(stderr, "Fast Reload failed (last app kept running): %v\n", err)
		return false
	}
	fmt.Fprintf(stdout, "Fast Reload complete (%s) in %s\n", platform, time.Since(started).Round(time.Millisecond))
	return true
}

func devCommand(root, platform string, reset bool, runner commandRunner, stdout, stderr io.Writer) error {
	if platform != "ios" && platform != "android" {
		return fmt.Errorf("unsupported dev platform %q", platform)
	}
	session, err := loadDevSession(root, reset)
	if err != nil {
		return err
	}
	previousSession, hadSession := os.LookupEnv("GONATIVE_RELOAD_SESSION")
	previousSkip, hadSkip := os.LookupEnv("GONATIVE_SKIP_FRAMEWORK_BUILD")
	defer restoreEnv("GONATIVE_RELOAD_SESSION", previousSession, hadSession)
	defer restoreEnv("GONATIVE_SKIP_FRAMEWORK_BUILD", previousSkip, hadSkip)
	_ = os.Setenv("GONATIVE_RELOAD_SESSION", session)
	if reset {
		fmt.Fprintln(stdout, "Fast Reload: starting with a fresh state session")
	}

	started := time.Now()
	if err := platformCommand(root, "run", platform, runner, stdout, stderr); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "Fast Reload ready (%s) in %s; watching for changes...\n", platform, time.Since(started).Round(time.Millisecond))
	_ = os.Setenv("GONATIVE_SKIP_FRAMEWORK_BUILD", "1")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return watchSources(ctx, root, devPollInterval, func(changed []string) {
		started := time.Now()
		fmt.Fprintf(stdout, "Fast Reload: %s changed; rebuilding...\n", strings.Join(changed, ", "))
		if err := platformCommand(root, "run", platform, runner, stdout, stderr); err != nil {
			fmt.Fprintf(stderr, "Fast Reload build failed (last app kept running): %v\n", err)
			return
		}
		fmt.Fprintf(stdout, "Fast Reload complete in %s\n", time.Since(started).Round(time.Millisecond))
	})
}

func watchSources(ctx context.Context, root string, interval time.Duration, changed func([]string)) error {
	previous, err := scanSources(root)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			next, err := scanSources(root)
			if err != nil {
				continue
			}
			changes := diffSnapshots(previous, next)
			if len(changes) > 0 {
				previous = next
				changed(changes)
			}
		}
	}
}

func scanSources(root string) (sourceSnapshot, error) {
	snapshot := make(sourceSnapshot)
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if rel != "." && ignoredDevDirectory(entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if !watchedDevFile(rel) {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		snapshot[filepath.ToSlash(rel)] = sha256.Sum256(data)
		return nil
	})
	return snapshot, err
}

func ignoredDevDirectory(name string) bool {
	switch name {
	case ".git", ".build", "build", ".gonative", ".gradle", "DerivedData", ".idea", ".vscode":
		return true
	default:
		return false
	}
}

func watchedDevFile(rel string) bool {
	base := filepath.Base(rel)
	if base == "go.mod" || base == "go.sum" || base == "gonative.yaml" {
		return true
	}
	if strings.HasSuffix(base, ".go") {
		return true
	}
	parts := strings.Split(filepath.ToSlash(rel), "/")
	return len(parts) > 1 && parts[0] == "assets"
}

func diffSnapshots(before, after sourceSnapshot) []string {
	set := make(map[string]bool)
	for path, hash := range after {
		if previous, ok := before[path]; !ok || previous != hash {
			set[path] = true
		}
	}
	for path := range before {
		if _, ok := after[path]; !ok {
			set[path] = true
		}
	}
	result := make([]string, 0, len(set))
	for path := range set {
		result = append(result, path)
	}
	sort.Strings(result)
	return result
}

func newDevSession() (string, error) {
	var bytes [12]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", fmt.Errorf("create Fast Reload session: %w", err)
	}
	return hex.EncodeToString(bytes[:]), nil
}

func loadDevSession(root string, reset bool) (string, error) {
	path := filepath.Join(root, ".gonative", "dev-session")
	if !reset {
		if data, err := os.ReadFile(path); err == nil {
			if session := strings.TrimSpace(string(data)); session != "" {
				return session, nil
			}
		}
	}
	session, err := newDevSession()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return "", fmt.Errorf("create Fast Reload metadata: %w", err)
	}
	if err := os.WriteFile(path, []byte(session+"\n"), 0o600); err != nil {
		return "", fmt.Errorf("write Fast Reload session: %w", err)
	}
	return session, nil
}

func restoreEnv(key, value string, existed bool) {
	if existed {
		_ = os.Setenv(key, value)
	} else {
		_ = os.Unsetenv(key)
	}
}
