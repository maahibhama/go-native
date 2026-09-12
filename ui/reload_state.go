package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ReloadStateOptions enables development-only state restoration. Applications
// normally configure this from generated debug bridge code.
type ReloadStateOptions struct {
	Path      string
	SessionID string
	OnWarning func(string)
}

type reloadSnapshot struct {
	Version   uint8                      `json:"version"`
	SessionID string                     `json:"sessionId"`
	Values    map[string]json.RawMessage `json:"values"`
}

type reloadStateStore struct {
	mu        sync.Mutex
	path      string
	sessionID string
	values    map[string]json.RawMessage
	warn      func(string)
	flush     *time.Timer
}

var activeReloadState struct {
	sync.RWMutex
	store *reloadStateStore
}

// ConfigureReloadState installs or disables the development reload-state
// store. An empty path or session ID disables persistence.
func ConfigureReloadState(options ReloadStateOptions) {
	if options.Path == "" || options.SessionID == "" {
		activeReloadState.Lock()
		activeReloadState.store = nil
		activeReloadState.Unlock()
		return
	}
	store := &reloadStateStore{path: options.Path, sessionID: options.SessionID, values: make(map[string]json.RawMessage), warn: options.OnWarning}
	store.load()
	activeReloadState.Lock()
	activeReloadState.store = store
	activeReloadState.Unlock()
}

func restoreReloadValue[T any](key string, initial T) (T, func(T)) {
	activeReloadState.RLock()
	store := activeReloadState.store
	activeReloadState.RUnlock()
	if store == nil {
		return initial, nil
	}
	store.mu.Lock()
	raw := append(json.RawMessage(nil), store.values[key]...)
	store.mu.Unlock()
	value := initial
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &value); err != nil {
			store.warning(fmt.Sprintf("reload state %q ignored: %v", key, err))
			value = initial
		}
	}
	return value, func(next T) { store.saveValue(key, next) }
}

func (s *reloadStateStore) load() {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		s.warning("reload state could not be read: " + err.Error())
		return
	}
	var snapshot reloadSnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil || snapshot.Version != 1 || snapshot.SessionID != s.sessionID {
		s.warning("reload state snapshot is incompatible and was ignored")
		return
	}
	s.values = snapshot.Values
	if s.values == nil {
		s.values = make(map[string]json.RawMessage)
	}
}

func (s *reloadStateStore) saveValue(key string, value any) {
	raw, err := json.Marshal(value)
	if err != nil {
		s.warning(fmt.Sprintf("reload state %q could not be encoded: %v", key, err))
		return
	}
	s.mu.Lock()
	s.values[key] = raw
	if s.flush != nil {
		s.flush.Stop()
	}
	s.flush = time.AfterFunc(25*time.Millisecond, s.flushValues)
	s.mu.Unlock()
}

func (s *reloadStateStore) flushValues() {
	s.mu.Lock()
	s.flush = nil
	snapshot := reloadSnapshot{Version: 1, SessionID: s.sessionID, Values: s.values}
	data, err := json.Marshal(snapshot)
	if err == nil {
		err = os.MkdirAll(filepath.Dir(s.path), 0o700)
	}
	if err == nil {
		temporary := s.path + ".tmp"
		err = os.WriteFile(temporary, data, 0o600)
		if err == nil {
			err = os.Rename(temporary, s.path)
		}
	}
	s.mu.Unlock()
	if err != nil {
		s.warning("reload state could not be saved: " + err.Error())
	}
}

func (s *reloadStateStore) warning(message string) {
	if s.warn != nil {
		s.warn(message)
	}
}
