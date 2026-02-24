# goxploit

**GoXploit** is a modular Go client for interacting with the [Metasploit RPC API](https://docs.rapid7.com/metasploit/standard-api-methods-reference).  
It includes:
- low-level RPC primitives (`pkg/msfrpc`)
- generated raw wrappers (`pkg/msfrpc/generated`)
- a typed SDK layer (`pkg/metasploit`)

## Features

- RPC client for Metasploit (`pkg/msfrpc`)
- Automatic method generation from Rapid7's API docs (`cmd/generator`)
- Typed SDK surface for stable app usage (`pkg/metasploit`)
- Easy integration with other Go projects
- Simple structure for extending or contributing

## Installation

```bash
go get github.com/arisinghackers/goxploit
```

## Usage 
Authenticate and call the typed SDK (`core.version`)
```go
import (
    "fmt"
    "log"

    "github.com/arisinghackers/goxploit/pkg/metasploit"
    "github.com/arisinghackers/goxploit/pkg/msfrpc"
)

client := msfrpc.NewMsfRpcClient("your_password", "false", "your_username", "127.0.0.1", 55552, "/api")
sdk := metasploit.NewClient(client)
login, err := sdk.Auth.Login("your_username", "your_password")
if err != nil {
    log.Fatalf("Auth error: %v", err)
}
fmt.Println("Token:", login.Token)

version, err := sdk.Core.Version()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Core Version: %s (ruby=%s api=%s)\n", version.Version, version.Ruby, version.API)
```
##### See /examples folder

## Project Structure

- `pkg/msfrpc`: low-level transport, auth, and generic request helpers.
- `pkg/msfrpc/generated`: generated map-based wrappers. Do not edit manually.
- `pkg/metasploit`: typed API layer intended for application code.
- `internal/generator`: generator internals (scraper + codegen logic).
- `cmd/generator`: CLI entrypoint for generation.

## Development Commands

```bash
make test
make generate
make check
```

## Contributing
Contributions are welcome. Please keep changes modular and follow idiomatic Go practices.
