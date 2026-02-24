package metasploit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arisinghackers/goxploit/pkg/msfrpc"
	"github.com/vmihailenco/msgpack/v5"
)

func TestAuthLoginSuccessSetsToken(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = msgpack.NewEncoder(w).Encode(map[string]interface{}{
			"result": "success",
			"token":  "abc123",
		})
	}))
	defer server.Close()

	rpc := msfrpc.NewMsfRpcClient("pw", "false", "user", "127.0.0.1", 55553, "/api/")
	rpc.BaseURL = server.URL
	rpc.HTTPClient = server.Client()

	client := NewClient(rpc)
	resp, err := client.Auth.Login("user", "pw")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if resp.Token != "abc123" {
		t.Fatalf("unexpected token: %q", resp.Token)
	}
	if rpc.GetToken() == nil || *rpc.GetToken() != "abc123" {
		t.Fatalf("expected token to be set on rpc client")
	}
}

func TestAuthLoginFailure(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = msgpack.NewEncoder(w).Encode(map[string]interface{}{
			"result":        "failure",
			"error_message": "bad credentials",
		})
	}))
	defer server.Close()

	rpc := msfrpc.NewMsfRpcClient("pw", "false", "user", "127.0.0.1", 55553, "/api/")
	rpc.BaseURL = server.URL
	rpc.HTTPClient = server.Client()

	client := NewClient(rpc)
	_, err := client.Auth.Login("user", "wrong")
	if err == nil || err.Error() != "bad credentials" {
		t.Fatalf("expected bad credentials error, got %v", err)
	}
}
