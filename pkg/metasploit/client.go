package metasploit

import "github.com/arisinghackers/goxploit/pkg/msfrpc"

// Client exposes typed Metasploit RPC services.
type Client struct {
	rpc  *msfrpc.MsfRpcClient
	Core *CoreService
}

func NewClient(rpc *msfrpc.MsfRpcClient) *Client {
	return &Client{
		rpc:  rpc,
		Core: &CoreService{rpc: rpc},
	}
}
