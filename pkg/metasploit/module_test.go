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

func TestModuleExecuteSuccess(t *testing.T) {
	t.Parallel()

	var payload []interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if err := msgpack.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		_ = msgpack.NewEncoder(w).Encode(map[string]interface{}{
			"result": "success",
			"job_id": int64(42),
			"uuid":   "run-uuid",
		})
	}))
	defer server.Close()

	rpc := msfrpc.NewMsfRpcClient("pw", "false", "user", "127.0.0.1", 55553, "/api/")
	rpc.BaseURL = server.URL
	rpc.HTTPClient = server.Client()
	rpc.SetToken("tok")

	client := NewClient(rpc)
	resp, err := client.Module.Execute(ExecuteModuleRequest{
		ModuleType: "exploit",
		ModuleName: "unix/ftp/vsftpd_234_backdoor",
		Options: map[string]interface{}{
			"RHOSTS": "127.0.0.1",
			"RPORT":  21,
		},
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if resp.JobID != 42 || resp.UUID != "run-uuid" {
		t.Fatalf("unexpected response: %+v", resp)
	}

	if len(payload) != 5 {
		t.Fatalf("unexpected payload length: %d", len(payload))
	}
	if asString(payload[0]) != "module.execute" {
		t.Fatalf("expected module.execute, got %v", payload[0])
	}
	if asString(payload[1]) != "tok" {
		t.Fatalf("expected token in second position, got %v", payload[1])
	}
}

func TestModuleExecuteFailure(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = msgpack.NewEncoder(w).Encode(map[string]interface{}{
			"result":        "failure",
			"error_message": "execution failed",
		})
	}))
	defer server.Close()

	rpc := msfrpc.NewMsfRpcClient("pw", "false", "user", "127.0.0.1", 55553, "/api/")
	rpc.BaseURL = server.URL
	rpc.HTTPClient = server.Client()
	rpc.SetToken("tok")

	client := NewClient(rpc)
	_, err := client.Module.Execute(ExecuteModuleRequest{
		ModuleType: "exploit",
		ModuleName: "bad",
	})
	if err == nil || err.Error() != "execution failed" {
		t.Fatalf("expected execution failed error, got %v", err)
	}
}

func TestModuleExecuteContextCancellation(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		_ = msgpack.NewEncoder(w).Encode(map[string]interface{}{
			"result": "success",
			"job_id": int64(1),
			"uuid":   "u",
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

	_, err := client.Module.ExecuteContext(ctx, ExecuteModuleRequest{
		ModuleType: "exploit",
		ModuleName: "x",
	})
	if err == nil || !strings.Contains(err.Error(), "context canceled") {
		t.Fatalf("expected context canceled error, got %v", err)
	}
}

func asString(v interface{}) string {
	switch typed := v.(type) {
	case string:
		return typed
	case []byte:
		return string(typed)
	default:
		return ""
	}
}
