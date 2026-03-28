package metasploit

import "github.com/arisinghackers/goxploit/pkg/msfrpc"

// Client exposes typed Metasploit RPC services.
type Client struct {
	rpc    *msfrpc.MsfRpcClient
	Auth   *AuthService
	Core   *CoreService
	Module *ModuleService
}

func NewClient(rpc *msfrpc.MsfRpcClient) *Client {
	return &Client{
		rpc:    rpc,
		Auth:   &AuthService{rpc: rpc},
		Core:   &CoreService{rpc: rpc},
		Module: &ModuleService{rpc: rpc},
	}
}
