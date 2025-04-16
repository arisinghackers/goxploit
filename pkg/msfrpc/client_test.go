package msfrpc_test

import (
	"github.com/arisinghackers/goxploit/pkg/msfrpc"
	"testing"
)

func TestAuthentication(t *testing.T) {
	client := msfrpc.NewMsfRpcClient("killodds", "false", "msf", "127.0.0.1", 55553, "/api/")
	token, err := client.MsfAuth()
	if err != nil {
		t.Fatal("Auth failed:", err)
	}
	if token == "" {
		t.Fatal("Empty token returned")
	}
}
