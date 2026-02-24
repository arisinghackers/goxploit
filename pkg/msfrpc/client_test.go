package msfrpc

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

func TestAuthenticatedRequestAddsTokenAfterMethod(t *testing.T) {
	t.Parallel()

	var received []interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if err := msgpack.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode request: %v", err)
		}

		_ = msgpack.NewEncoder(w).Encode(map[string]interface{}{"result": "success"})
	}))
	defer server.Close()

	client := NewMsfRpcClient("pw", "false", "user", "127.0.0.1", 55553, "/api/")
	client.BaseURL = server.URL
	client.HTTPClient = server.Client()
	client.SetToken("token-123")

	_, err := client.AuthenticatedRequest([]any{"core.version", "arg1"})
	if err != nil {
		t.Fatalf("AuthenticatedRequest failed: %v", err)
	}

	if len(received) != 3 {
		t.Fatalf("expected 3 request args, got %d", len(received))
	}

	if got := asString(received[0]); got != "core.version" {
		t.Fatalf("expected method first, got %q", got)
	}
	if got := asString(received[1]); got != "token-123" {
		t.Fatalf("expected token second, got %q", got)
	}
	if got := asString(received[2]); got != "arg1" {
		t.Fatalf("expected arg third, got %q", got)
	}
}

func TestMsfAuthParsesStringAndBytesToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		response map[string]interface{}
		expected string
	}{
		{name: "string token", response: map[string]interface{}{"token": "abc"}, expected: "abc"},
		{name: "bytes token", response: map[string]interface{}{"token": []byte("xyz")}, expected: "xyz"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = msgpack.NewEncoder(w).Encode(tc.response)
			}))
			defer server.Close()

			client := NewMsfRpcClient("pw", "false", "user", "127.0.0.1", 55553, "/api/")
			client.BaseURL = server.URL
			client.HTTPClient = server.Client()

			token, err := client.MsfAuth()
			if err != nil {
				t.Fatalf("MsfAuth failed: %v", err)
			}
			if token != tc.expected {
				t.Fatalf("expected token %q, got %q", tc.expected, token)
			}
		})
	}
}

func TestMsfRequestReturnsHTTPStatusError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server exploded", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewMsfRpcClient("pw", "false", "user", "127.0.0.1", 55553, "/api/")
	client.BaseURL = server.URL
	client.HTTPClient = server.Client()

	_, err := client.MsfRequest([]interface{}{"core.version"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "status 500") {
		t.Fatalf("expected status error, got %v", err)
	}
}

func TestAuthenticatedRequestValidatesInput(t *testing.T) {
	t.Parallel()

	client := NewMsfRpcClient("pw", "false", "user", "127.0.0.1", 55553, "/api/")

	_, err := client.AuthenticatedRequest([]any{"core.version"})
	if err == nil || !strings.Contains(err.Error(), "token is nil") {
		t.Fatalf("expected missing token error, got %v", err)
	}

	client.SetToken("tok")
	_, err = client.AuthenticatedRequest([]any{})
	if err == nil || !strings.Contains(err.Error(), "payload cannot be empty") {
		t.Fatalf("expected empty payload error, got %v", err)
	}
}

func asString(v interface{}) string {
	switch value := v.(type) {
	case string:
		return value
	case []byte:
		return string(value)
	default:
		return fmt.Sprintf("%v", value)
	}
}

func TestMsfRequestDecodesMsgpackBody(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload []interface{}
		if err := msgpack.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if asString(payload[0]) != "core.version" {
			t.Fatalf("unexpected method: %v", payload[0])
		}

		var buf bytes.Buffer
		if err := msgpack.NewEncoder(&buf).Encode(map[string]interface{}{"version": "6.4.0"}); err != nil {
			t.Fatalf("encode response: %v", err)
		}
		_, _ = w.Write(buf.Bytes())
	}))
	defer server.Close()

	client := NewMsfRpcClient("pw", "false", "user", "127.0.0.1", 55553, "/api/")
	client.BaseURL = server.URL
	client.HTTPClient = server.Client()

	resp, err := client.MsfRequest([]interface{}{"core.version"})
	if err != nil {
		t.Fatalf("MsfRequest failed: %v", err)
	}
	if got := asString(resp["version"]); got != "6.4.0" {
		t.Fatalf("expected decoded version, got %q", got)
	}
}

func TestMsfRequestContextCancellation(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		_ = msgpack.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
	}))
	defer server.Close()

	client := NewMsfRpcClient("pw", "false", "user", "127.0.0.1", 55553, "/api/")
	client.BaseURL = server.URL
	client.HTTPClient = server.Client()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.MsfRequestContext(ctx, []interface{}{"core.version"})
	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}
	if !strings.Contains(err.Error(), "context canceled") {
		t.Fatalf("expected context canceled error, got %v", err)
	}
}
