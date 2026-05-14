package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/arisinghackers/goxploit/pkg/metasploit"
	"github.com/arisinghackers/goxploit/pkg/msfrpc"
)

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	host := env("MSF_HOST", "127.0.0.1")
	port := 3000
	ssl := env("MSF_SSL", "false")
	webURI := env("MSF_WEB_URI", "/api")
	user := env("MSF_USER", "msf")
	pass := env("MSF_PASS", "msf")

	client := msfrpc.NewMsfRpcClient(pass, ssl, user, host, port, webURI)
	typedClient := metasploit.NewClient(client)
	ctx := context.Background()

	login, err := typedClient.Auth.LoginContext(ctx, user, pass)
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
