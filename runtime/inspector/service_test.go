package inspector

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	gonativeruntime "github.com/go-native/go-native/runtime"
	"github.com/go-native/go-native/ui"
)

type fakeSource struct{}

type closedListener struct{ address net.Addr }

func (l closedListener) Accept() (net.Conn, error) { return nil, errors.New("closed") }
func (l closedListener) Close() error              { return nil }
func (l closedListener) Addr() net.Addr            { return l.address }

func (fakeSource) LogEntries() []gonativeruntime.LogEntry {
	return []gonativeruntime.LogEntry{{Kind: gonativeruntime.LogBatchApplied, Sequence: 4, MutationCount: 2}}
}

func (fakeSource) TreeSnapshot() gonativeruntime.TreeSnapshot {
	return gonativeruntime.SnapshotTree(&ui.Node{ID: 3, Type: ui.NodeText, Props: ui.Props{Text: "Hello"}})
}

func TestServiceServesReadOnlyDiagnostics(t *testing.T) {
	service := New(fakeSource{}, "")
	response := httptest.NewRecorder()
	service.routes().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/v1/tree", nil))
	result := response.Result()
	defer result.Body.Close()
	if result.Header.Get("Cache-Control") != "no-store" {
		t.Fatalf("cache policy = %q", result.Header.Get("Cache-Control"))
	}
	var tree gonativeruntime.TreeSnapshot
	if err := json.NewDecoder(result.Body).Decode(&tree); err != nil {
		t.Fatal(err)
	}
	if tree.Root == nil || tree.Root.Props.Text != "Hello" {
		t.Fatalf("unexpected tree: %#v", tree)
	}

	response = httptest.NewRecorder()
	service.routes().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/v1/logs", nil))
	var logs struct {
		Entries []gonativeruntime.LogEntry `json:"entries"`
	}
	if err := json.NewDecoder(response.Body).Decode(&logs); err != nil {
		t.Fatal(err)
	}
	if len(logs.Entries) != 1 || logs.Entries[0].Sequence != 4 {
		t.Fatalf("unexpected logs: %#v", logs.Entries)
	}
}

func TestServiceRejectsNonLoopbackAndInvalidLifecycle(t *testing.T) {
	for _, address := range []string{"0.0.0.0:0", "[::]:0", "example.com:9000", "bad-address"} {
		if _, err := New(fakeSource{}, address).Start(); err == nil {
			t.Errorf("Start(%q) succeeded", address)
		}
	}
	if _, err := New(nil, "").Start(); err == nil {
		t.Fatal("nil source started")
	}

	originalListen := listen
	listen = func(_, _ string) (net.Listener, error) {
		return closedListener{address: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 49152}}, nil
	}
	t.Cleanup(func() { listen = originalListen })
	service := New(fakeSource{}, "127.0.0.1:0")
	address, err := service.Start()
	if err != nil || address != "127.0.0.1:49152" {
		t.Fatalf("Start = %q, %v", address, err)
	}
	if _, err := service.Start(); err == nil {
		t.Fatal("second Start succeeded")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := service.Stop(ctx); err != nil {
		t.Fatal(err)
	}
	if err := service.Stop(ctx); err != nil {
		t.Fatalf("second Stop: %v", err)
	}
}

func TestServiceRejectsMutationMethods(t *testing.T) {
	service := New(fakeSource{}, "")
	request := httptest.NewRequest(http.MethodPost, "/v1/logs", nil)
	response := httptest.NewRecorder()
	service.routes().ServeHTTP(response, request)
	_, _ = io.Copy(io.Discard, response.Result().Body)
	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST status = %d, want 405", response.Code)
	}
}
