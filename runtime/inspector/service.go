// Package inspector exposes opt-in, read-only runtime diagnostics over loopback HTTP.
package inspector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	gonativeruntime "github.com/go-native/go-native/runtime"
)

// Source is the read-only runtime surface consumed by the inspector.
type Source interface {
	LogEntries() []gonativeruntime.LogEntry
	TreeSnapshot() gonativeruntime.TreeSnapshot
}

// Service owns an opt-in loopback HTTP inspector.
type Service struct {
	source Source
	addr   string

	mu       sync.Mutex
	listener net.Listener
	server   *http.Server
}

var listen = net.Listen

// New creates a stopped inspector. An empty address selects 127.0.0.1 with an
// ephemeral port. Non-loopback addresses are rejected by Start.
func New(source Source, address string) *Service {
	if address == "" {
		address = "127.0.0.1:0"
	}
	return &Service{source: source, addr: address}
}

// Start begins serving and returns the selected loopback address.
func (s *Service) Start() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.source == nil {
		return "", errors.New("inspector: nil diagnostics source")
	}
	if s.listener != nil {
		return "", errors.New("inspector: already started")
	}
	if err := validateLoopbackAddress(s.addr); err != nil {
		return "", err
	}
	listener, err := listen("tcp", s.addr)
	if err != nil {
		return "", fmt.Errorf("inspector: listen: %w", err)
	}
	server := &http.Server{
		Handler:           s.routes(),
		ReadHeaderTimeout: 2 * time.Second,
		IdleTimeout:       15 * time.Second,
		MaxHeaderBytes:    16 << 10,
	}
	s.listener = listener
	s.server = server
	go func() { _ = server.Serve(listener) }()
	return listener.Addr().String(), nil
}

func (s *Service) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/tree", s.serveTree)
	mux.HandleFunc("GET /v1/logs", s.serveLogs)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	return mux
}

// Stop gracefully closes the inspector. It is safe to call more than once.
func (s *Service) Stop(ctx context.Context) error {
	s.mu.Lock()
	server := s.server
	if server == nil {
		s.mu.Unlock()
		return nil
	}
	s.server = nil
	s.listener = nil
	s.mu.Unlock()
	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("inspector: shutdown: %w", err)
	}
	return nil
}

func (s *Service) serveTree(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, s.source.TreeSnapshot())
}

func (s *Service) serveLogs(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, struct {
		Entries []gonativeruntime.LogEntry `json:"entries"`
	}{Entries: s.source.LogEntries()})
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(true)
	if err := encoder.Encode(value); err != nil {
		http.Error(w, "failed to encode diagnostics", http.StatusInternalServerError)
	}
}

func validateLoopbackAddress(address string) error {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("inspector: invalid address %q: %w", address, err)
	}
	if strings.EqualFold(host, "localhost") {
		return nil
	}
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		return fmt.Errorf("inspector: address %q is not loopback", address)
	}
	return nil
}
