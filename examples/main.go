package main

// running some tests on the lib
import (
	"context"
	"fmt"
	"log"

	"github.com/arisinghackers/goxploit/pkg/metasploit"
	"github.com/arisinghackers/goxploit/pkg/msfrpc"
)

func main() {

	client := msfrpc.NewMsfRpcClient("killodds", "false", "veemweaver", "127.0.0.1", 3000, "/api")
	typedClient := metasploit.NewClient(client)
	ctx := context.Background()

	login, err := typedClient.Auth.LoginContext(ctx, "veemweaver", "killodds")
	if err != nil {
		log.Fatalf("Error authenticating: %v", err)
	}
	log.Printf("Authenticated successfully, token: %s", login.Token)

	version, err := typedClient.Core.VersionContext(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Metasploit Core Version: %s (ruby=%s api=%s)\n", version.Version, version.Ruby, version.API)

}
