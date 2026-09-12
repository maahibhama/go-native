package main

import (
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
)

const devPollInterval = 200 * time.Millisecond

type sourceSnapshot map[string][32]byte

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
