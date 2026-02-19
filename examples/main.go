package main

// running some tests on the lib
import (
	"fmt"
	"log"

	"github.com/arisinghackers/goxploit/pkg/metasploit"
	"github.com/arisinghackers/goxploit/pkg/msfrpc"
)

func main() {

	client := msfrpc.NewMsfRpcClient("killodds", "false", "veemweaver", "127.0.0.1", 3000, "/api")
	resp, err := client.MsfAuth()
	if err != nil {
		log.Fatalf("Error authenticating: %v", err)
	}
	log.Printf("Authenticated successfully, token: %s", resp)

	typedClient := metasploit.NewClient(client)
	version, err := typedClient.Core.Version()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Metasploit Core Version: %s (ruby=%s api=%s)\n", version.Version, version.Ruby, version.API)

}
