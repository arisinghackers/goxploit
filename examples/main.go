package main

// running some tests on the lib
import (
	"fmt"
	"log"

	"github.com/arisinghackers/goxploit/pkg/msfrpc"
)

func main() {

	client := msfrpc.NewMsfRpcClient("killodds", "false", "veemweaver", "127.0.0.1", 3000, "/api")
	resp, err := client.MsfAuth()
	if err != nil {
		log.Fatalf("Error authenticating: %v", err)
	}
	log.Printf("Authenticated successfully, token: %s", resp)

	response, err := client.AuthenticatedRequest([]any{"core.version"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Metasploit Core Version: %+v\n", response)

	scraper := msfrpc.NewMsfPayloadScraper()
	methods, err := scraper.GetArraysPayloadsFromWebsite()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(methods)

}
