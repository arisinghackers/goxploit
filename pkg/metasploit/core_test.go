package metasploit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/arisinghackers/goxploit/pkg/msfrpc"
	"github.com/vmihailenco/msgpack/v5"
)

func TestCoreVersionSuccess(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = msgpack.NewEncoder(w).Encode(map[string]interface{}{
			"version": "6.4.0-dev",
			"ruby":    "3.1.5",
			"api":     "1.0",
		})
	}))
	defer server.Close()

	rpc := msfrpc.NewMsfRpcClient("pw", "false", "user", "127.0.0.1", 55553, "/api/")
	rpc.BaseURL = server.URL
	rpc.HTTPClient = server.Client()
	rpc.SetToken("tok")

	client := NewClient(rpc)
	got, err := client.Core.Version()
	if err != nil {
		t.Fatalf("Version failed: %v", err)
	}

	if got.Version != "6.4.0-dev" || got.Ruby != "3.1.5" || got.API != "1.0" {
		t.Fatalf("unexpected typed result: %+v", got)
	}
}

func TestCoreVersionMissingRequiredField(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = msgpack.NewEncoder(w).Encode(map[string]interface{}{
			"ruby": "3.1.5",
		})
	}))
	defer server.Close()

	rpc := msfrpc.NewMsfRpcClient("pw", "false", "user", "127.0.0.1", 55553, "/api/")
	rpc.BaseURL = server.URL
	rpc.HTTPClient = server.Client()
	rpc.SetToken("tok")

	client := NewClient(rpc)
	_, err := client.Core.Version()
	if err == nil || !strings.Contains(err.Error(), "missing \"version\"") {
		t.Fatalf("expected missing version error, got %v", err)
	}
}

func TestCoreVersionInvalidType(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = msgpack.NewEncoder(w).Encode(map[string]interface{}{
			"version": 123,
		})
	}))
	defer server.Close()

	rpc := msfrpc.NewMsfRpcClient("pw", "false", "user", "127.0.0.1", 55553, "/api/")
	rpc.BaseURL = server.URL
	rpc.HTTPClient = server.Client()
	rpc.SetToken("tok")

	client := NewClient(rpc)
	_, err := client.Core.Version()
	if err == nil || !strings.Contains(err.Error(), "invalid type") {
		t.Fatalf("expected invalid type error, got %v", err)
	}
}

func TestCoreVersionContextCancellation(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		_ = msgpack.NewEncoder(w).Encode(map[string]interface{}{
			"version": "6.4.0-dev",
		})
	}))
	defer server.Close()

	rpc := msfrpc.NewMsfRpcClient("pw", "false", "user", "127.0.0.1", 55553, "/api/")
	rpc.BaseURL = server.URL
	rpc.HTTPClient = server.Client()
	rpc.SetToken("tok")

	client := NewClient(rpc)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.Core.VersionContext(ctx)
	if err == nil || !strings.Contains(err.Error(), "context canceled") {
		t.Fatalf("expected context canceled error, got %v", err)
	}
}
