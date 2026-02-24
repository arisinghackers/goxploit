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
	typedClient := metasploit.NewClient(client)
	login, err := typedClient.Auth.Login("veemweaver", "killodds")
	if err != nil {
		log.Fatalf("Error authenticating: %v", err)
	}
	log.Printf("Authenticated successfully, token: %s", login.Token)

	version, err := typedClient.Core.Version()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Metasploit Core Version: %s (ruby=%s api=%s)\n", version.Version, version.Ruby, version.API)

}
